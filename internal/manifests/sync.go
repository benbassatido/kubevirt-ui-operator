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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	consoleuiv1alpha1 "github.com/kubevirt-ui/kubevirt-ui-operator/api/v1alpha1"
)

// ReconcileStack applies Helm-parity operands for the effective spec in Context.
func ReconcileStack(ctx context.Context, c client.Client, scheme *runtime.Scheme, m *Context) error {
	if pc := PriorityClass(m); pc != nil {
		if err := reconcileCluster(ctx, c, func() client.Object { return PriorityClass(m) }); err != nil {
			return err
		}
	}

	steps := []func() error{
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return pluginServiceAccount(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return proxyServiceAccount(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return nginxConfigMap(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return userSettingsConfigMap(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return featuresConfigMap(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return pluginDeployment(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return proxyDeployment(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return pluginService(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return proxyService(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return configReaderRole(m) })
		},
		func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return configReaderRoleBinding(m) })
		},
	}
	if consoleuiv1alpha1.BoolOrDefault(m.Spec.NetworkPolicy.Enabled, true) {
		steps = append(steps,
			func() error {
				return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return pluginNetworkPolicy(m) })
			},
			func() error {
				return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return proxyNetworkPolicy(m) })
			},
		)
	} else {
		steps = append(steps,
			func() error {
				return deleteNamespaced(ctx, c, m.OperandNamespace, "kubevirt-console-plugin-np", &networkingv1.NetworkPolicy{})
			},
			func() error {
				return deleteNamespaced(ctx, c, m.OperandNamespace, "kubevirt-apiserver-proxy-np", &networkingv1.NetworkPolicy{})
			},
		)
	}
	if consoleuiv1alpha1.BoolOrDefault(m.Spec.Route.Enabled, true) {
		steps = append(steps, func() error {
			return reconcileNamespaced(ctx, c, scheme, m, func() client.Object { return proxyRoute(m) })
		})
	} else {
		steps = append(steps, func() error {
			return deleteNamespaced(ctx, c, m.OperandNamespace, m.Spec.Route.Name, &routev1.Route{})
		})
	}

	steps = append(steps, func() error {
		return reconcileCluster(ctx, c, func() client.Object { return consolePlugin(m) })
	})

	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
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

func deleteNamespaced(ctx context.Context, c client.Client, ns, name string, obj client.Object) error {
	obj.SetNamespace(ns)
	obj.SetName(name)
	err := c.Delete(ctx, obj)
	if errors.IsNotFound(err) {
		return nil
	}
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
