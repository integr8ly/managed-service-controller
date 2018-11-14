package v1alpha1

import (
	appsv1 "github.com/openshift/api/apps/v1"
	authv1 "github.com/openshift/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

//Removed as the managed-service-controller needs to grant these permissions
//var fuseServiceRole = &rbacv1beta1.Role{
//	ObjectMeta: metav1.ObjectMeta{
//		Name: "syndesis-operator",
//		Labels: map[string]string{
//			"app":                   "syndesis",
//			"syndesis.io/app":       "syndesis",
//			"syndesis.io/type":      "operator",
//			"syndesis.io/component": "syndesis-operator",
//		},
//	},
//	Rules: []rbacv1beta1.PolicyRule{
//		{
//			APIGroups: []string{"syndesis.io"},
//			Resources: []string{"syndesises", "syndesises/finalizers"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{""},
//			Resources: []string{"pods", "services", "endpoints", "persistentvolumeclaims", "configmaps", "secrets", "serviceaccounts"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{""},
//			Resources: []string{"events"},
//			Verbs:     []string{"get", "list"},
//		},
//		{
//			APIGroups: []string{"rbac.authorization.k8s.io"},
//			Resources: []string{"rolebindings"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"template.openshift.io"},
//			Resources: []string{"processedtemplates"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"image.openshift.io"},
//			Resources: []string{"imagestreams"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"apps.openshift.io"},
//			Resources: []string{"deploymentconfigs"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"build.openshift.io"},
//			Resources: []string{"buildconfigs"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"authorization.openshift.io"},
//			Resources: []string{"rolebindings"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//		{
//			APIGroups: []string{"route.openshift.io"},
//			Resources: []string{"routes", "routes/custom-host"},
//			Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
//		},
//	},
//}

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
