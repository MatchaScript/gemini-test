package operator

import (
	"context"
	"fmt"
	"math/rand"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	nanokubev1alpha1 "github.com/MatchaScript/nanokube/internal/apis/nanokube.io/v1alpha1"
)

// NanokubeMachineReconciler reconciles a NanokubeMachine object by claiming an matching InventoryMachine.
type NanokubeMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *NanokubeMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var nm nanokubev1alpha1.NanokubeMachine
	if err := r.Get(ctx, req.NamespacedName, &nm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Already claimed
	if nm.Status.InventoryMachineRef != "" {
		return ctrl.Result{}, nil
	}

	// Find available InventoryMachine matching hostSelector
	var invList nanokubev1alpha1.InventoryMachineList
	if err := r.List(ctx, &invList); err != nil {
		return ctrl.Result{}, fmt.Errorf("list inventory machines: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(&nm.Spec.HostSelector)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("invalid hostSelector: %w", err)
	}

	var candidates []nanokubev1alpha1.InventoryMachine
	for _, inv := range invList.Items {
		if inv.Status.MachineRef == "" && !inv.Spec.Maintenance && selector.Matches(labels.Set(inv.Labels)) {
			candidates = append(candidates, inv)
		}
	}

	if len(candidates) == 0 {
		logger.Info("no available InventoryMachine matching hostSelector", "selector", nm.Spec.HostSelector)
		return ctrl.Result{Requeue: true}, nil
	}

	// Select a random candidate to claim (Metal3 pattern)
	chosen := candidates[rand.Intn(len(candidates))]

	// Claim
	chosen.Status.MachineRef = nm.Name
	if err := r.Status().Update(ctx, &chosen); err != nil {
		return ctrl.Result{}, fmt.Errorf("update inventory machine status: %w", err)
	}

	nm.Spec.ProviderID = chosen.Spec.ProviderID
	if nm.Spec.ProviderID == "" {
		nm.Spec.ProviderID = fmt.Sprintf("nanokube://%s", chosen.Spec.Address)
	}

	if err := r.Update(ctx, &nm); err != nil {
		return ctrl.Result{}, fmt.Errorf("update nanokube machine: %w", err)
	}

	nm.Status.InventoryMachineRef = chosen.Name
	nm.Status.Initialization.Provisioned = true
	nm.Status.Ready = true

	if err := r.Status().Update(ctx, &nm); err != nil {
		return ctrl.Result{}, fmt.Errorf("update nanokube machine status: %w", err)
	}

	logger.Info("claimed InventoryMachine for NanokubeMachine", "inventoryMachine", chosen.Name, "nanokubeMachine", nm.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NanokubeMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nanokubev1alpha1.NanokubeMachine{}).
		Complete(r)
}

// BootstrapSecretName derives the deterministic bootstrap secret name for a machine.
func BootstrapSecretName(machineName string) string {
	return fmt.Sprintf("%s-bootstrap-data", machineName)
}

// BuildBootstrapSecret creates a per-Machine Secret carrying rendered desired bytes + revision.
func BuildBootstrapSecret(namespace, machineName string, desiredBlob []byte, revision string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BootstrapSecretName(machineName),
			Namespace: namespace,
			Labels: map[string]string{
				"nanokube.io/revision": revision,
			},
		},
		Data: map[string][]byte{
			"value":    desiredBlob,
			"revision": []byte(revision),
		},
	}
}
