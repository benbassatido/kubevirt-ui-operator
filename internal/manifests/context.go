package manifests

import "github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests/builders"

// Context holds operand manifest inputs after spec defaults (see builders package).
type Context = builders.Context

// NewContext deep-copies the CR and applies default spec values for manifest generation.
var NewContext = builders.NewContext
