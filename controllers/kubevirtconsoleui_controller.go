package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	routev1 "github.com/openshift/api/route/v1"

	consolev1 "github.com/openshift/api/console/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/operconfig"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/status"
)

// KubeVirtConsoleUIReconciler reconciles the console UI stack for a KubeVirtConsoleUI CR.
type KubeVirtConsoleUIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Config operconfig.Config
}

// +kubebuilder:rbac:groups=consoleui.kubevirt.io,resources=kubevirtconsoleuis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=consoleui.kubevirt.io,resources=kubevirtconsoleuis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=consoleui.kubevirt.io,resources=kubevirtconsoleuis/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=serviceaccounts;services;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scheduling.k8s.io,resources=priorityclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=console.openshift.io,resources=consoleplugins,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.openshift.io,resources=consoles,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch

func (r *KubeVirtConsoleUIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if req.Namespace != r.Config.OperandNamespace {
		return ctrl.Result{}, nil
	}

	cr := &consoleuiv1alpha1.KubeVirtConsoleUI{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	mctx := manifests.NewContext(cr)
	reconcileErr := r.reconcile(ctx, cr, mctx)
	deploymentsWaiting := ""
	if reconcileErr == nil {
		deploymentsWaiting = status.OperandDeploymentsWaiting(ctx, r.Client,
			mctx.OperandNamespace, mctx.Spec.Plugin.Name, mctx.Spec.Proxy.Name)
	}
	return r.updateStatusAndRequeue(ctx, cr, reconcileErr, deploymentsWaiting)
}

func (r *KubeVirtConsoleUIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	operandNS := r.Config.OperandNamespace
	inOperandNS := builder.WithPredicates(predicate.NewPredicateFuncs(func(o client.Object) bool {
		return o.GetNamespace() == operandNS
	}))

	return ctrl.NewControllerManagedBy(mgr).
		For(&consoleuiv1alpha1.KubeVirtConsoleUI{}, inOperandNS).
		Owns(&appsv1.Deployment{}, inOperandNS).
		Owns(&corev1.Service{}, inOperandNS).
		Owns(&corev1.ConfigMap{}, inOperandNS).
		Owns(&corev1.ServiceAccount{}, inOperandNS).
		Owns(&rbacv1.Role{}, inOperandNS).
		Owns(&rbacv1.RoleBinding{}, inOperandNS).
		Owns(&networkingv1.NetworkPolicy{}, inOperandNS).
		Owns(&routev1.Route{}, inOperandNS).
		Watches(
			&consolev1.ConsolePlugin{},
			handler.EnqueueRequestsFromMapFunc(r.mapConsolePluginToRequests),
		).
		Watches(
			&schedulingv1.PriorityClass{},
			handler.EnqueueRequestsFromMapFunc(r.mapPriorityClassToRequests),
		).
		Watches(
			&operatorv1.Console{},
			handler.EnqueueRequestsFromMapFunc(r.mapConsoleToRequests),
		).
		Complete(r)
}
