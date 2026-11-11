package hcoversion

import "testing"

func TestOperatorVersionFromStatus(t *testing.T) {
	hc := &HyperConverged{
		Status: HyperConvergedStatus{
			Versions: []ComponentVersion{
				{Name: "kubevirt", Version: "1.2.3"},
				{Name: OperatorVersionName, Version: "4.20.0"},
			},
		},
	}
	if got := operatorVersionFromStatus(hc); got != "4.20.0" {
		t.Fatalf("got %q", got)
	}
}

func TestOperatorVersionFromStatus_missing(t *testing.T) {
	if got := operatorVersionFromStatus(&HyperConverged{}); got != "" {
		t.Fatalf("got %q", got)
	}
}
