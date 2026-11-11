package manifests

import (
	"context"

	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/features"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/hcoversion"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/manifests/builders"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

// ReconcileStack applies Helm-parity operands for the effective spec in Context.
// hcoReader should be a direct API reader (not cache-scoped) for HyperConverged discovery.
func ReconcileStack(ctx context.Context, c client.Client, hcoReader client.Reader, scheme *runtime.Scheme, m *Context) error {
	if err := reconcileCluster(ctx, c, func() client.Object { return builders.PriorityClass(m) }); err != nil {
		return err
	}

	steps := []func() error{
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.PluginServiceAccount(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ProxyServiceAccount(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.NginxConfigMap(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.UserSettingsConfigMap(m) })
		},
		func() error {
			return reconcileFeaturesConfigMap(ctx, c, hcoReader, scheme, m)
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.PluginDeployment(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ProxyDeployment(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.PluginService(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ProxyService(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ConfigReaderRole(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ConfigReaderRoleBinding(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.PluginNetworkPolicy(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ProxyNetworkPolicy(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return builders.ProxyRoute(m) })
		},
		func() error {
			return reconcileCluster(ctx, c, func() client.Object { return builders.ConsolePlugin(m) })
		},
	}

	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func reconcileFeaturesConfigMap(ctx context.Context, c client.Client, hcoReader client.Reader, scheme *runtime.Scheme, m *Context) error {
	hcoReader = coalesceReader(hcoReader, c)
	hcoVersion := ""
	manageHCOVersion := false
	if m.HCONamespace != "" {
		installed, err := hcoversion.Installed(ctx, hcoReader)
		if err != nil {
			return err
		}
		if installed {
			manageHCOVersion = true
			hcoVersion, err = hcoversion.OperatorVersion(ctx, hcoReader, m.HCONamespace)
			if err != nil {
				manageHCOVersion = false
			}
		}
	}
	cm := &corev1.ConfigMap{}
	cm.Name = names.FeaturesConfigMap
	cm.Namespace = m.OperandNamespace
	_, err := controllerutil.CreateOrUpdate(ctx, c, cm, func() error {
		cm.Labels = builders.FeaturesConfigMapShell(m, nil).Labels
		cm.Data = features.MergeData(cm.Data, hcoVersion, manageHCOVersion)
		return controllerutil.SetControllerReference(m.CR, cm, scheme)
	})
	return err
}

func coalesceReader(preferred, fallback client.Reader) client.Reader {
	if preferred != nil {
		return preferred
	}
	return fallback
}

func reconcileNamespaced(ctx context.Context, c client.Client, scheme *runtime.Scheme, m *Context, build func() client.Object) error {
	obj := build()
	_, err := controllerutil.CreateOrUpdate(ctx, c, obj, func() error {
		fresh := build()
		obj.SetLabels(fresh.GetLabels())
		obj.SetAnnotations(fresh.GetAnnotations())
		copySpec(obj, fresh)
		return controllerutil.SetControllerReference(m.CR, obj, scheme)
	})
	return err
}

func reconcileCluster(ctx context.Context, c client.Client, build func() client.Object) error {
	obj := build()
	_, err := controllerutil.CreateOrUpdate(ctx, c, obj, func() error {
		fresh := build()
		obj.SetLabels(fresh.GetLabels())
		obj.SetAnnotations(fresh.GetAnnotations())
		copySpec(obj, fresh)
		return nil
	})
	return err
}

func copySpec(dst, src client.Object) {
	switch d := dst.(type) {
	case *corev1.ServiceAccount:
		s := src.(*corev1.ServiceAccount)
		d.Secrets = s.Secrets
		d.ImagePullSecrets = s.ImagePullSecrets
	case *corev1.ConfigMap:
		s := src.(*corev1.ConfigMap)
		d.Data = s.Data
		d.BinaryData = s.BinaryData
	case *appsv1.Deployment:
		s := src.(*appsv1.Deployment)
		d.Spec = s.Spec
	case *corev1.Service:
		s := src.(*corev1.Service)
		d.Spec = s.Spec
	case *rbacv1.Role:
		s := src.(*rbacv1.Role)
		d.Rules = s.Rules
	case *rbacv1.RoleBinding:
		s := src.(*rbacv1.RoleBinding)
		d.RoleRef = s.RoleRef
		d.Subjects = s.Subjects
	case *networkingv1.NetworkPolicy:
		s := src.(*networkingv1.NetworkPolicy)
		d.Spec = s.Spec
	case *routev1.Route:
		s := src.(*routev1.Route)
		d.Spec = s.Spec
	case *schedulingv1.PriorityClass:
		s := src.(*schedulingv1.PriorityClass)
		d.Value = s.Value
		d.GlobalDefault = s.GlobalDefault
		d.Description = s.Description
		d.PreemptionPolicy = s.PreemptionPolicy
	case *consolev1.ConsolePlugin:
		s := src.(*consolev1.ConsolePlugin)
		d.Spec = s.Spec
	}
}
