package builders

import (
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
)

func ConfigReaderRole(m *Context) *rbacv1.Role {
	return &rbacv1.Role{
		Name:      m.Spec.RBAC.ConfigReaderRoleName,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{m.Spec.Config.UserSettingsConfigMapName},
				Verbs:         []string{"get", "update", "patch"},
			},
			{
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{m.Spec.Config.FeaturesConfigMapName},
				Verbs:         []string{"get"},
			},
		},
	}
}

func ConfigReaderRoleBinding(m *Context) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		Name:      m.Spec.RBAC.ConfigReaderRoleBindingName,
		Namespace: m.OperandNamespace,
		Labels:    labels.Plugin(m.CR),
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     m.Spec.RBAC.ConfigReaderRoleName,
		},
		Subjects: []rbacv1.Subject{{
			APIGroup: rbacv1.GroupName,
			Kind:     "Group",
			Name:     "system:authenticated",
		}},
	}
}
