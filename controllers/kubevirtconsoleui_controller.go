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
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/hcoversion"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
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
// +kubebuilder:rbac:groups=hco.kubevirt.io,resources=hyperconvergeds,verbs=get;list;watch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch

func (r *KubeVirtConsoleUIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if req.Name != names.ResourceName {
		return ctrl.Result{}, nil
	}

	cr := &consoleuiv1alpha1.KubeVirtConsoleUI{}
	if err := r.Get(ctx, client.ObjectKey{Name: req.Name}, cr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	mctx := manifests.NewContext(cr, r.Config)
	reconcileErr := r.reconcile(ctx, cr, mctx)
	deploymentsWaiting := ""
	if reconcileErr == nil {
		deploymentsWaiting = status.OperandDeploymentsWaiting(ctx, r.Client,
			mctx.OperandNamespace, names.PluginDeployment, names.ProxyDeployment)
	}
	return r.updateStatusAndRequeue(ctx, cr, reconcileErr, deploymentsWaiting)
}

func (r *KubeVirtConsoleUIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	operandNS := r.Config.OperandNamespace
	inOperandNS := predicate.NewPredicateFuncs(func(o client.Object) bool {
		return o.GetNamespace() == operandNS
	})
	clusterCR := predicate.NewPredicateFuncs(func(o client.Object) bool {
		return o.GetName() == names.ResourceName
	})

	ctrlBuilder := ctrl.NewControllerManagedBy(mgr).
		For(&consoleuiv1alpha1.KubeVirtConsoleUI{}, builder.WithPredicates(clusterCR)).
		Owns(&appsv1.Deployment{}, builder.WithPredicates(inOperandNS)).
		Owns(&corev1.Service{}, builder.WithPredicates(inOperandNS)).
		Owns(&corev1.ConfigMap{}, builder.WithPredicates(inOperandNS)).
		Owns(&corev1.ServiceAccount{}, builder.WithPredicates(inOperandNS)).
		Owns(&rbacv1.Role{}, builder.WithPredicates(inOperandNS)).
		Owns(&rbacv1.RoleBinding{}, builder.WithPredicates(inOperandNS)).
		Owns(&networkingv1.NetworkPolicy{}, builder.WithPredicates(inOperandNS)).
		Owns(&routev1.Route{}, builder.WithPredicates(inOperandNS)).
		Watches(
			&consolev1.ConsolePlugin{},
			handler.EnqueueRequestsFromMapFunc(r.mapConsolePluginToRequests),
			builder.WithPredicates(predicate.NewPredicateFuncs(func(o client.Object) bool {
				return o.GetName() == names.ConsolePluginName
			})),
		).
		Watches(
			&schedulingv1.PriorityClass{},
			handler.EnqueueRequestsFromMapFunc(r.mapPriorityClassToRequests),
			builder.WithPredicates(predicate.NewPredicateFuncs(func(o client.Object) bool {
				return o.GetName() == names.PriorityClass
			})),
		).
		Watches(
			&operatorv1.Console{},
			handler.EnqueueRequestsFromMapFunc(r.mapConsoleToRequests),
			builder.WithPredicates(predicate.NewPredicateFuncs(func(o client.Object) bool {
				return o.GetName() == "cluster"
			})),
		).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.mapFeaturesConfigMapToRequests),
			builder.WithPredicates(predicate.NewPredicateFuncs(func(o client.Object) bool {
				return o.GetNamespace() == operandNS && o.GetName() == names.FeaturesConfigMap
			})),
		)

	if r.Config.HCOInstalled && r.Config.HCONamespace != "" {
		inHCONS := predicate.NewPredicateFuncs(func(o client.Object) bool {
			return o.GetNamespace() == r.Config.HCONamespace
		})
		ctrlBuilder = ctrlBuilder.Watches(
			&hcoversion.HyperConverged{},
			handler.EnqueueRequestsFromMapFunc(r.mapHyperConvergedToRequests),
			builder.WithPredicates(inHCONS),
		)
	}

	return ctrlBuilder.Complete(r)
}
