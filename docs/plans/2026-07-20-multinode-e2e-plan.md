# Plan: Multi-Node E2E Test Harness & Production Readiness

## Overview
Implement a multi-node E2E test harness inspired by `bink` (bootc-dev/bink) and `bootloose` (k0sproject/bootloose).
This harness runs multiple containerized bootc/systemd nodes representing 1 Control Plane (`cp-1`) and 2 Worker Nodes (`worker-1`, `worker-2`), executing the full single-path pipeline (`nanokube init` -> `nanokube add-node` -> `push` -> `confext refresh` -> pod scheduling) and validating cluster resilience.

## Test Topology & Lifecycle
1. **Infrastructure Provisioning**:
   - Spin up containerized nodes with `systemd` enabled and `nanokube-agent` running over gRPC.
2. **Control Plane Initialization**:
   - Execute `nanokube init` targeting `cp-1`.
   - Verify control plane static pods (`apiserver`, `etcd`, `controller-manager`, `scheduler`) and kubelet.
3. **Worker Joining**:
   - Generate bootstrap token via `nanokube` on `cp-1`.
   - Execute `nanokube add-node` targeting `worker-1` and `worker-2`.
   - Verify node registration in Kubernetes cluster (`kubectl get nodes`).
4. **Workload & In-Place Update Verification**:
   - Deploy sample workload deployment to multi-node cluster.
   - Execute `push.DesiredToAgent` live update to verify `systemd-confext refresh --mutable=yes` across control plane and worker nodes without node teardown.
   - Verify non-leakage of secrets and disaster recovery procedures.

## Task Breakdown
- [x] Task 1: Create multi-node harness definitions and helper routines in `test/e2e/scenarios_multinode_test.go`.
- [x] Task 2: Implement multi-node E2E scenario suite (`scenarios_multinode_test.go`).
- [x] Task 3: Add `multinode-e2e` coverage to `.github/workflows/ci.yaml`.
- [ ] Task 4: Verify full test suite locally and on GitHub Actions until 100% ALL GREEN.
