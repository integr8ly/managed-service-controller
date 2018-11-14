package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	appsv1 "github.com/openshift/api/apps/v1"
)
var integrationServiceAccount = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Name: "integration-controller",
	},
}

//var integrationServiceRole = &rbacv1beta1.Role{
//	ObjectMeta: metav1.ObjectMeta{
//		Name: "integration-controller",
//	},
//	Rules: []rbacv1beta1.PolicyRule{
//		{
//			APIGroups: []string{"integreatly.org"},
//			Resources: []string{"*"},
//			Verbs:     []string{"*"},
//		},
//		{
//			APIGroups: []string{""},
//			Resources: []string{"pods", "services", "endpoints", "persistentvolumeclaims", "configmaps", "secrets"},
//			Verbs:     []string{"*"},
//		},
//		{
//			APIGroups: []string{"apps"},
//			Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
//			Verbs:     []string{"*"},
//		},
//		{
//			APIGroups: []string{"syndesis.io"},
//			Resources: []string{"*"},
//			Verbs:     []string{"*"},
//		},
//	},
//}

var integrationServiceRoleBinding = &rbacv1beta1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Name: "ns-integration-controller",
	},
	Subjects: []rbacv1beta1.Subject{
		{
			Kind: "ServiceAccount",
			Name: "integration-controller",
		},
	},
	RoleRef: rbacv1beta1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "integration-controller",
		APIGroup: "rbac.authorization.k8s.io",
	},
}

var integrationDeploymentConfig = &appsv1.DeploymentConfig{
	ObjectMeta: metav1.ObjectMeta{
		Name: IntegrationControllerName,
	},
	Spec: appsv1.DeploymentConfigSpec{
		Strategy: appsv1.DeploymentStrategy{
			Type: "Recreate",
		},
		Replicas: 1,
		Selector: map[string]string{
			"name": IntegrationControllerName,
		},
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name": IntegrationControllerName,
				},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: IntegrationControllerName,
				Containers: []corev1.Container{
					{
						Name: "integration-controller",
						Image: "quay.io/integreatly/" + IntegrationControllerName + ":dev",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 60000,
								Name: "metrics",
							},
						},
						Command: []string{
							IntegrationControllerName,
							"--allow-insecure=true",
							"--log-level=debug",
						},
						ImagePullPolicy: "Always",
						Env: []corev1.EnvVar{
							{
								Name: "WATCH_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name: "OPERATOR_NAME",
								Value: IntegrationControllerName,
							},
							{
								Name: "USER_NAMESPACES",
								Value: "",
							},
						},
					},
				},
			},
		},
	},
}