package defaults

import (
	"testing"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"k8s.io/utils/ptr"
)

func TestApply_defaultTrueBools(t *testing.T) {
	spec := &consoleuiv1alpha1.KubeVirtConsoleUISpec{}
	Apply(spec)

	if !ptr.Deref(spec.HCOLabels.Enabled, false) {
		t.Fatal("expected hcoLabels.enabled true by default")
	}
	if !ptr.Deref(spec.PriorityClass.Create, false) {
		t.Fatal("expected priorityClass.create true by default")
	}
	if !ptr.Deref(spec.Console.RegisterPlugin, false) {
		t.Fatal("expected console.registerPlugin true by default")
	}
	if !ptr.Deref(spec.NetworkPolicy.Enabled, false) {
		t.Fatal("expected networkPolicy.enabled true by default")
	}
	if !ptr.Deref(spec.Route.Enabled, false) {
		t.Fatal("expected route.enabled true by default")
	}
}

func TestApply_explicitFalse(t *testing.T) {
	spec := &consoleuiv1alpha1.KubeVirtConsoleUISpec{
		NetworkPolicy: consoleuiv1alpha1.NetworkPolicySpec{
			Enabled: new(false),
		},
	}
	Apply(spec)
	if ptr.Deref(spec.NetworkPolicy.Enabled, true) {
		t.Fatal("expected networkPolicy.enabled to stay false")
	}
}

func TestApply_defaultImages(t *testing.T) {
	spec := &consoleuiv1alpha1.KubeVirtConsoleUISpec{}
	Apply(spec)
	if spec.Plugin.Image != defaultPluginImage {
		t.Fatalf("plugin image: got %q want %q", spec.Plugin.Image, defaultPluginImage)
	}
	if spec.Proxy.Image != defaultProxyImage {
		t.Fatalf("proxy image: got %q want %q", spec.Proxy.Image, defaultProxyImage)
	}
}

func TestApply_explicitZeroReplicas(t *testing.T) {
	spec := &consoleuiv1alpha1.KubeVirtConsoleUISpec{
		Plugin: consoleuiv1alpha1.PluginSpec{
			Replicas: new(int32(0)),
		},
		Proxy: consoleuiv1alpha1.ProxySpec{
			Replicas: new(int32(0)),
		},
		PriorityClass: consoleuiv1alpha1.PriorityClassSpec{
			Value: new(int32(0)),
		},
	}
	Apply(spec)
	if ptr.Deref(spec.Plugin.Replicas, 1) != 0 {
		t.Fatal("expected plugin replicas to stay 0")
	}
	if ptr.Deref(spec.Proxy.Replicas, 1) != 0 {
		t.Fatal("expected proxy replicas to stay 0")
	}
	if ptr.Deref(spec.PriorityClass.Value, 1000000000) != 0 {
		t.Fatal("expected priorityClass value to stay 0")
	}
}
