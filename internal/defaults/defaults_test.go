package defaults

import (
	"testing"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
)

func TestApply_noPanicOnEmptySpec(t *testing.T) {
	spec := &consoleuiv1alpha1.KubeVirtConsoleUISpec{}
	Apply(spec)
}
