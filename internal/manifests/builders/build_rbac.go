package builders

import (
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/labels"
	"github.com/kubevirt-ui/kubevirt-ui-operator/internal/names"
)

func ConfigReaderRole(m *Context) *rbacv1.Role {
	return &rbacv1.Role{
		Name:      names.ConfigReaderRole,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:         []string{"get"},
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{names.UserSettingsConfigMap},
			},
			{
				Verbs:         []string{"get"},
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{names.FeaturesConfigMap},
			},
		},
	}
}

func ConfigReaderRoleBinding(m *Context) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		Name:      names.ConfigReaderRoleBinding,
		Namespace: m.OperandNamespace,
		Labels:    labels.UIConfig(m.CR),
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     names.ConfigReaderRole,
		},
		Subjects: []rbacv1.Subject{{
			Kind:     rbacv1.GroupKind,
			APIGroup: rbacv1.GroupName,
			Name:     "system:authenticated",
		}},
	}
}
