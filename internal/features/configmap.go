package features

import "maps"

// ConfigMap data keys for kubevirt-ui-features (Helm parity).
const (
	KeyAutomaticSubscriptionActivationKey  = "automaticSubscriptionActivationKey"
	KeyAutomaticSubscriptionOrganizationId = "automaticSubscriptionOrganizationId"
	KeyDisabledGuestSystemLogsAccess       = "disabledGuestSystemLogsAccess"
	KeyKubevirtApiserverProxy              = "kubevirtApiserverProxy"
	KeyLoadBalancerEnabled                 = "loadBalancerEnabled"
	KeyNodePortAddress                     = "nodePortAddress"
	KeyNodePortEnabled                     = "nodePortEnabled"
	KeyIPStackType                         = "ipStackType"
	KeyHCOVersion                          = "hcoVersion"
)

// DefaultData returns initial feature ConfigMap entries when the ConfigMap is first created.
func DefaultData() map[string]string {
	return map[string]string{
		KeyAutomaticSubscriptionActivationKey:  "",
		KeyAutomaticSubscriptionOrganizationId: "",
		KeyDisabledGuestSystemLogsAccess:       "false",
		KeyKubevirtApiserverProxy:              "true",
		KeyLoadBalancerEnabled:                 "true",
		KeyNodePortAddress:                     "",
		KeyNodePortEnabled:                     "false",
	}
}

// MergeData preserves existing ConfigMap data and fills missing default keys.
// When manageHCOVersion is true, sets or removes hcoVersion from the merged data.
// When manageHCOVersion is false, an existing hcoVersion entry is left unchanged.
func MergeData(existing map[string]string, hcoVersion string, manageHCOVersion bool) map[string]string {
	out := maps.Clone(existing)
	if out == nil {
		out = make(map[string]string, len(DefaultData()))
	}
	for k, v := range DefaultData() {
		if _, ok := out[k]; !ok {
			out[k] = v
		}
	}
	if !manageHCOVersion {
		return out
	}
	if hcoVersion != "" {
		out[KeyHCOVersion] = hcoVersion
	} else {
		delete(out, KeyHCOVersion)
	}
	return out
}
