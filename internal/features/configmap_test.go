package features

import "testing"

func TestMergeData_preservesUserEdits(t *testing.T) {
	existing := map[string]string{
		KeyKubevirtApiserverProxy: "false",
		KeyNodePortEnabled:        "true",
	}
	merged := MergeData(existing, "")
	if merged[KeyKubevirtApiserverProxy] != "false" {
		t.Fatalf("expected user value preserved, got %q", merged[KeyKubevirtApiserverProxy])
	}
	if merged[KeyDisabledGuestSystemLogsAccess] != "false" {
		t.Fatal("expected default key filled")
	}
}

func TestMergeData_clusterHCOVersion(t *testing.T) {
	existing := map[string]string{KeyHCOVersion: "1.0"}
	merged := MergeData(existing, "2.0")
	if merged[KeyHCOVersion] != "2.0" {
		t.Fatalf("got %q", merged[KeyHCOVersion])
	}
}

func TestMergeData_createUsesDefaults(t *testing.T) {
	merged := MergeData(nil, "")
	if merged[KeyLoadBalancerEnabled] != "true" {
		t.Fatal("expected defaults on create")
	}
	if _, ok := merged[KeyHCOVersion]; ok {
		t.Fatal("expected hcoVersion omitted when HCO version is unknown")
	}
}

func TestMergeData_omitsHCOVersionWhenUndetermined(t *testing.T) {
	existing := map[string]string{KeyHCOVersion: "1.0"}
	merged := MergeData(existing, "")
	if _, ok := merged[KeyHCOVersion]; ok {
		t.Fatalf("expected hcoVersion omitted, got %q", merged[KeyHCOVersion])
	}
}
