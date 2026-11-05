package defaults

import consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"

// Apply fills unset spec fields. No-op today; kept for manifest context parity.
func Apply(spec *consoleuiv1alpha1.KubeVirtConsoleUISpec) {}
