package operator

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	nanokubev1alpha1 "github.com/MatchaScript/nanokube/internal/apis/nanokube.io/v1alpha1"
)

func TestNanokubeMachineReconciler_ClaimsMatchingInventoryMachine(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = nanokubev1alpha1.AddToScheme(scheme)

	inv := &nanokubev1alpha1.InventoryMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-1",
			Labels: map[string]string{
				"nanokube.io/pool": "worker",
			},
		},
		Spec: nanokubev1alpha1.InventoryMachineSpec{
			Address: "192.168.1.100",
		},
	}

	nm := &nanokubev1alpha1.NanokubeMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "machine-1",
			Namespace: "default",
		},
		Spec: nanokubev1alpha1.NanokubeMachineSpec{
			HostSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"nanokube.io/pool": "worker",
				},
			},
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(inv, nm).WithObjects(inv, nm).Build()

	r := &NanokubeMachineReconciler{
		Client: client,
		Scheme: scheme,
	}

	ctx := context.Background()
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "machine-1", Namespace: "default"}}

	_, err := r.Reconcile(ctx, req)
	if err != nil {
		t.Fatalf("Reconcile error: %v", err)
	}

	var updatedNM nanokubev1alpha1.NanokubeMachine
	if err := client.Get(ctx, req.NamespacedName, &updatedNM); err != nil {
		t.Fatalf("Get NanokubeMachine: %v", err)
	}

	if updatedNM.Status.InventoryMachineRef != "node-1" {
		t.Errorf("InventoryMachineRef = %q, want node-1", updatedNM.Status.InventoryMachineRef)
	}
	if !updatedNM.Status.Ready {
		t.Errorf("Ready = false, want true")
	}
}

func TestBuildBootstrapSecret_FormatAndData(t *testing.T) {
	secret := BuildBootstrapSecret("default", "machine-1", []byte("blob-data"), "rev-123")

	if secret.Name != "machine-1-bootstrap-data" {
		t.Errorf("Secret.Name = %q, want machine-1-bootstrap-data", secret.Name)
	}
	if string(secret.Data["value"]) != "blob-data" {
		t.Errorf("Secret.Data[value] = %q, want blob-data", string(secret.Data["value"]))
	}
	if string(secret.Data["revision"]) != "rev-123" {
		t.Errorf("Secret.Data[revision] = %q, want rev-123", string(secret.Data["revision"]))
	}
}
