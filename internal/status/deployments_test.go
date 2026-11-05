package status

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestDeploymentAvailable(t *testing.T) {
	replicas := int32(2)
	t.Run("available condition true", func(t *testing.T) {
		ok, _ := deploymentAvailable(&appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{Replicas: &replicas},
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: 2,
				Conditions: []appsv1.DeploymentCondition{{
					Type:   appsv1.DeploymentAvailable,
					Status: corev1.ConditionTrue,
				}},
			},
		})
		if !ok {
			t.Fatal("expected available")
		}
	})
	t.Run("insufficient replicas", func(t *testing.T) {
		ok, detail := deploymentAvailable(&appsv1.Deployment{
			Name: "kubevirt-console-plugin",
			Spec: appsv1.DeploymentSpec{Replicas: &replicas},
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: 0,
				Conditions: []appsv1.DeploymentCondition{{
					Type:    appsv1.DeploymentAvailable,
					Status:  corev1.ConditionFalse,
					Reason:  "MinimumReplicasUnavailable",
					Message: "Deployment does not have minimum availability.",
				}},
			},
		})
		if ok {
			t.Fatal("expected not available")
		}
		if detail == "" {
			t.Fatal("expected detail message")
		}
	})
}
