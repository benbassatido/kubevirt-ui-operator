package hcoversion

import (
	"context"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInstalled_crdPresent(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apiextensionsv1.AddToScheme(scheme)

	crd := &apiextensionsv1.CustomResourceDefinition{
		Name: CRDName,
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd).Build()

	installed, err := Installed(context.Background(), c)
	if err != nil || !installed {
		t.Fatalf("expected HCO installed when CRD exists, err=%v installed=%v", err, installed)
	}
}

func TestInstalled_crdMissing(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = apiextensionsv1.AddToScheme(scheme)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	installed, err := Installed(context.Background(), c)
	if err != nil || installed {
		t.Fatalf("expected HCO not installed when CRD is absent, err=%v installed=%v", err, installed)
	}
}
