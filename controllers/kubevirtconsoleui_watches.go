package controllers

import (
	"context"

	consolev1 "github.com/openshift/api/console/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/hcoversion"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

func (r *KubeVirtConsoleUIReconciler) mapConsolePluginToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	if cp, ok := obj.(*consolev1.ConsolePlugin); ok && cp.GetName() == names.ConsolePluginName {
		return r.enqueueClusterCR()
	}
	return nil
}

func (r *KubeVirtConsoleUIReconciler) mapPriorityClassToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	if pc, ok := obj.(*schedulingv1.PriorityClass); ok && pc.GetName() == names.PriorityClass {
		return r.enqueueClusterCR()
	}
	return nil
}

func (r *KubeVirtConsoleUIReconciler) mapConsoleToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	if console, ok := obj.(*operatorv1.Console); ok && console.GetName() == "cluster" {
		return r.enqueueClusterCR()
	}
	return nil
}

func (r *KubeVirtConsoleUIReconciler) mapFeaturesConfigMapToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok || cm.GetNamespace() != r.Config.OperandNamespace || cm.GetName() != names.FeaturesConfigMap {
		return nil
	}
	return r.enqueueClusterCR()
}

func (r *KubeVirtConsoleUIReconciler) mapHyperConvergedCRDToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	crd, ok := obj.(*apiextensionsv1.CustomResourceDefinition)
	if !ok || crd.GetName() != hcoversion.CRDName {
		return nil
	}
	return r.enqueueClusterCR()
}

func (r *KubeVirtConsoleUIReconciler) mapHyperConvergedToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	if _, ok := obj.(*hcoversion.HyperConverged); ok {
		return r.enqueueClusterCR()
	}
	return nil
}

func (r *KubeVirtConsoleUIReconciler) enqueueClusterCR() []reconcile.Request {
	return []reconcile.Request{{Name: names.ResourceName}}
}
