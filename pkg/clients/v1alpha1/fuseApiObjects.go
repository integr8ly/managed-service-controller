package v1alpha1

import (
	appsv1 "github.com/openshift/api/apps/v1"
	authv1 "github.com/openshift/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getFuseServiceAccount(name, namespace string) *corev1.ServiceAccount {
	operatorName := name + "-operator"
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":                   name,
				"syndesis.io/app":       name,
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
		},
	}
}

func getFuseServiceRoleBinding(name, namespace string) *rbacv1beta1.RoleBinding {
	operatorName := name + "-operator"
	return &rbacv1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorName + ":install",
			Namespace: namespace,
			Labels: map[string]string{
				"app":                   name,
				"syndesis.io/app":       name,
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind: "ServiceAccount",
				Name: operatorName,
			},
		},
		RoleRef: rbacv1beta1.RoleRef{
			Kind:     "ClusterRole",
			Name:     operatorName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func getViewRoleBinding(name string) *authv1.RoleBinding {
	operatorName := name + "-operator"
	return &authv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName + ":view",
			Labels: map[string]string{
				"app":                   name,
				"syndesis.io/app":       name,
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind: "ServiceAccount",
				Name: operatorName,
			},
		},
		RoleRef: corev1.ObjectReference{
			Name: "view",
		},
	}
}

func getEditRoleBinding(name string) *authv1.RoleBinding {
	operatorName := name + "-operator"
	return &authv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName + ":edit",
			Labels: map[string]string{
				"app":                   name,
				"syndesis.io/app":       name,
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind: "ServiceAccount",
				Name: operatorName,
			},
		},
		RoleRef: corev1.ObjectReference{
			Name: "edit",
		},
	}
}

func getFuseDeploymentConfig(cfg map[string]string) *appsv1.DeploymentConfig {
	operatorName := cfg["name"] + "-operator"
	return &appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName,
			Labels: map[string]string{
				"app":                   cfg["name"],
				"syndesis.io/app":       cfg["name"],
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
		},
		Spec: appsv1.DeploymentConfigSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
			Replicas: 1,
			Selector: map[string]string{
				"syndesis.io/app":       cfg["name"],
				"syndesis.io/type":      "operator",
				"syndesis.io/component": operatorName,
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"syndesis.io/app":       cfg["name"],
						"syndesis.io/type":      "operator",
						"syndesis.io/component": operatorName,
					},
				},

				Spec: corev1.PodSpec{
					ServiceAccountName: operatorName,
					Containers: []corev1.Container{
						{
							Name:            operatorName,
							Image:           " ",
							ImagePullPolicy: "IfNotPresent",
							Env: []corev1.EnvVar{
								{
									Name: "WATCH_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
				},
			},
			Triggers: appsv1.DeploymentTriggerPolicies{
				appsv1.DeploymentTriggerPolicy{
					ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							operatorName,
						},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      cfg["imageName"] + ":" + cfg["imageTag"],
							Namespace: cfg["imageStreamNamespace"],
						},
					},
					Type: "ImageChange",
				},
				appsv1.DeploymentTriggerPolicy{
					Type: "ConfigChange",
				},
			},
		},
	}
}
