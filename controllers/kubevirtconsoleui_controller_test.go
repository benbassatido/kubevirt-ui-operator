package controllers

import (
	"context"
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/operconfig"
)

func TestReconcilePlatformRequirementsNotMetWithoutOpenShiftConsole(t *testing.T) {
	ctx := t.Context()
	ns := "reconcile-test"
	t.Setenv(operconfig.EnvOperandNamespace, ns)

	if err := k8sClient.Create(ctx, &corev1.Namespace{Name: ns}); err != nil {
		t.Fatalf("create namespace: %v", err)
	}
	t.Cleanup(func() {
		// t.Context() is canceled before Cleanup runs; use a live context for teardown.
		if err := k8sClient.Delete(context.Background(), &corev1.Namespace{Name: ns}); err != nil && !apierrors.IsNotFound(err) {
			t.Errorf("delete namespace %q: %v", ns, err)
		}
		_ = os.Unsetenv(operconfig.EnvOperandNamespace)
	})

	cr := &consoleuiv1alpha1.KubeVirtConsoleUI{
		Name: names.ResourceName,
	}
	if err := k8sClient.Create(ctx, cr); err != nil {
		t.Fatalf("create CR: %v", err)
	}

	r := &KubeVirtConsoleUIReconciler{
		Client:    k8sClient,
		APIReader: k8sClient,
		Scheme:    testScheme,
		Config:    operconfig.LoadFromEnv(),
	}
	if _, err := r.Reconcile(ctx, ctrl.Request{Name: names.ResourceName}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	updated := &consoleuiv1alpha1.KubeVirtConsoleUI{}
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: names.ResourceName}, updated); err != nil {
		t.Fatalf("get CR: %v", err)
	}

	var platformCond *metav1.Condition
	for i := range updated.Status.Conditions {
		if updated.Status.Conditions[i].Type == consoleuiv1alpha1.ConditionTypePlatformRequirementsMet {
			platformCond = &updated.Status.Conditions[i]
			break
		}
	}
	if platformCond == nil {
		t.Fatal("expected PlatformRequirementsMet condition")
	}
	if platformCond.Status != metav1.ConditionFalse {
		t.Fatalf("expected PlatformRequirementsMet=False, got %s", platformCond.Status)
	}
	if platformCond.Reason != "ConsolePluginCRDMissing" {
		t.Fatalf("expected ConsolePluginCRDMissing reason, got %q", platformCond.Reason)
	}
}
