package manifests

import (
	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/defaults"
)

// Context holds a deep copy of the CR after defaults.Apply (Helm parity). Use CR for labels and owner references; Spec mirrors CR.Spec.
type Context struct {
	CR               *consoleuiv1alpha1.KubeVirtConsoleUI
	Spec             consoleuiv1alpha1.KubeVirtConsoleUISpec
	OperandNamespace string
}

// NewContext deep-copies the CR and applies default spec values for manifest generation.
func NewContext(cr *consoleuiv1alpha1.KubeVirtConsoleUI) *Context {
	copy := cr.DeepCopy()
	defaults.Apply(&copy.Spec)
	return &Context{
		CR:               copy,
		Spec:             copy.Spec,
		OperandNamespace: copy.Namespace,
	}
}
