package v1alpha1

import (
	appsv1 "github.com/openshift/api/apps/v1"
	authv1 "github.com/openshift/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FUSE_IMAGE_STREAMS_NAMESPACE string = "openshift"

var fuseServiceAccount = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Name: "syndesis-operator",
		Labels: map[string]string{
			"app":                   "syndesis",
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
	},
}

var fuseServiceRoleBinding = &rbacv1beta1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Name: "syndesis-operator:install",
		Labels: map[string]string{
			"app":                   "syndesis",
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
	},
	Subjects: []rbacv1beta1.Subject{
		{
			Kind: "ServiceAccount",
			Name: "syndesis-operator",
		},
	},
	RoleRef: rbacv1beta1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "syndesis-operator",
		APIGroup: "rbac.authorization.k8s.io",
	},
}

var viewRoleBinding = &authv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Name: "syndesis-operator:view",
		Labels: map[string]string{
			"app":                   "syndesis",
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
	},
	Subjects: []corev1.ObjectReference{
		{
			Kind: "ServiceAccount",
			Name: "syndesis-operator",
		},
	},
	RoleRef: corev1.ObjectReference{
		Name: "view",
	},
}

var editRoleBinding = &authv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Name: "syndesis-operator:edit",
		Labels: map[string]string{
			"app":                   "syndesis",
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
	},
	Subjects: []corev1.ObjectReference{
		{
			Kind: "ServiceAccount",
			Name: "syndesis-operator",
		},
	},
	RoleRef: corev1.ObjectReference{
		Name: "edit",
	},
}

var fuseDeploymentConfig = &appsv1.DeploymentConfig{
	ObjectMeta: metav1.ObjectMeta{
		Name: "syndesis-operator",
		Labels: map[string]string{
			"app":                   "syndesis",
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
	},
	Spec: appsv1.DeploymentConfigSpec{
		Strategy: appsv1.DeploymentStrategy{
			Type: "Recreate",
		},
		Replicas: 1,
		Selector: map[string]string{
			"syndesis.io/app":       "syndesis",
			"syndesis.io/type":      "operator",
			"syndesis.io/component": "syndesis-operator",
		},
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"syndesis.io/app":       "syndesis",
					"syndesis.io/type":      "operator",
					"syndesis.io/component": "syndesis-operator",
				},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: "syndesis-operator",
				Containers: []corev1.Container{
					{
						Name:            "syndesis-operator",
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
						"syndesis-operator",
					},
					From: corev1.ObjectReference{
						Kind:      "ImageStreamTag",
						Name:      "fuse-online-operator:1.4",
						Namespace: FUSE_IMAGE_STREAMS_NAMESPACE,
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
