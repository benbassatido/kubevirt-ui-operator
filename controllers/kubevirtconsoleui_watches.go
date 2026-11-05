package controllers

import (
	"context"

	consolev1 "github.com/openshift/api/console/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests"
	"k8s.io/utils/ptr"
)

func (r *KubeVirtConsoleUIReconciler) mapConsolePluginToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	cp, ok := obj.(*consolev1.ConsolePlugin)
	if !ok {
		return nil
	}
	return r.requestsForConsolePluginName(ctx, cp.GetName())
}

func (r *KubeVirtConsoleUIReconciler) mapPriorityClassToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	pc, ok := obj.(*schedulingv1.PriorityClass)
	if !ok {
		return nil
	}
	return r.requestsForPriorityClassName(ctx, pc.GetName())
}

func (r *KubeVirtConsoleUIReconciler) mapConsoleToRequests(ctx context.Context, obj client.Object) []reconcile.Request {
	console, ok := obj.(*operatorv1.Console)
	if !ok || console.GetName() != "cluster" {
		return nil
	}
	return r.requestsForConsoleRegistration(ctx)
}

func (r *KubeVirtConsoleUIReconciler) requestsForConsolePluginName(ctx context.Context, pluginName string) []reconcile.Request {
	list, err := r.listConsoleUIs(ctx)
	if err != nil {
		return nil
	}
	var out []reconcile.Request
	for i := range list.Items {
		cr := &list.Items[i]
		mctx := manifests.NewContext(cr)
		if mctx.Spec.ConsolePlugin.Name == pluginName {
			out = append(out, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cr)})
		}
	}
	return out
}

func (r *KubeVirtConsoleUIReconciler) requestsForPriorityClassName(ctx context.Context, priorityClassName string) []reconcile.Request {
	list, err := r.listConsoleUIs(ctx)
	if err != nil {
		return nil
	}
	var out []reconcile.Request
	for i := range list.Items {
		cr := &list.Items[i]
		mctx := manifests.NewContext(cr)
		if !ptr.Deref(mctx.Spec.PriorityClass.Create, true) {
			continue
		}
		if mctx.Spec.PriorityClass.Name == priorityClassName {
			out = append(out, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cr)})
		}
	}
	return out
}

func (r *KubeVirtConsoleUIReconciler) requestsForConsoleRegistration(ctx context.Context) []reconcile.Request {
	list, err := r.listConsoleUIs(ctx)
	if err != nil {
		return nil
	}
	var out []reconcile.Request
	for i := range list.Items {
		cr := &list.Items[i]
		mctx := manifests.NewContext(cr)
		if !ptr.Deref(mctx.Spec.Console.RegisterPlugin, true) {
			continue
		}
		out = append(out, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cr)})
	}
	return out
}

func (r *KubeVirtConsoleUIReconciler) listConsoleUIs(ctx context.Context) (consoleuiv1alpha1.KubeVirtConsoleUIList, error) {
	var list consoleuiv1alpha1.KubeVirtConsoleUIList
	if err := r.List(ctx, &list); err != nil {
		log.FromContext(ctx).Error(err, "list KubeVirtConsoleUI for secondary watch")
		return list, err
	}
	return list, nil
}
