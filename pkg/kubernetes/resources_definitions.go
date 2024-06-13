package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// jobで環境変数を入れる時に使う構造体
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type RequestData struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Port    string   `json:"port"`
	EnvVars []EnvVar `json:"envVars"`
}

// deploymentのリソース設定
func DeploymentDefinition(appName string, deploymentName string, registryName string, envVars []EnvVar) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: registryName,
							Env:   EnvDefinition(envVars),
						},
					},
				},
			},
		},
	}
	return deployment
}

// serviceのリソース設定
func ServiceDefinition(appName string, serviceName string, targetPort int) *apiv1.Service {

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(targetPort),
				},
			},
			Selector: map[string]string{
				"app": appName,
			},
		},
	}
	return service
}
func EnvDefinition(envVars []EnvVar) []apiv1.EnvVar {
	var k8sEnvVars []apiv1.EnvVar
	for _, envVar := range envVars {
		k8sEnvVars = append(k8sEnvVars, apiv1.EnvVar{
			Name:  envVar.Name,
			Value: envVar.Value,
		})
	}
	return k8sEnvVars
}

// ingressのリソース設定
func IngressDefinition(ingressName string, hostName string, serviceName string) *networkingv1.Ingress {

	nginxServiceName := "nginx"
	pathType := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &nginxServiceName,
			Rules: []networkingv1.IngressRule{
				{
					Host: hostName,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}

// jobtのリソース設定
func JobDefinition(githubUrl string, appName string, registryName string, envVars []EnvVar) *batchv1.Job {

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko",
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "kaniko",
				},

				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{
						{
							Name:  "init-gitclone",
							Image: "alpine/git:2.43.4",
							Command: []string{
								"git",
								"clone",
								githubUrl,
								"/workspace",
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "dockerfile-storage",
									MountPath: "/workspace",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  "kaniko",
							Image: "gcr.io/kaniko-project/executor:latest",
							Args: []string{
								"--dockerfile=/workspace/Dockerfile",
								"--context=dir:///workspace",
								// "--no-push",
								// "--destination=docker.io/kouxi00/test:latest",
								"--destination=" + registryName + ":latest",
							},
							Env: EnvDefinition(envVars),
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "dockerfile-storage",
									MountPath: "/workspace",
								},
								{
									Name:      "kaniko-secret",
									MountPath: "/kaniko/.docker",
								},
							},
						},
					},
					RestartPolicy: "Never",
					Volumes: []apiv1.Volume{
						{
							Name: "dockerfile-storage",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: appName + "-pvc",
								},
							},
						},
						{
							Name: "kaniko-secret",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "dockerhub-secret",
									Items: []apiv1.KeyToPath{
										{
											Key:  ".dockerconfigjson",
											Path: "config.json",
										},
									},
								},
							},
						},
					},
				},
			},
			BackoffLimit:            int32Ptr(0),
			TTLSecondsAfterFinished: int32Ptr(20),
		},
	}
	return job
}

// pvcのリソース設定
func PvcDefinition(Name string, Uid string, appName string) *apiv1.PersistentVolumeClaim {
	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-pvc",
			Annotations: map[string]string{
				"volume.kubernetes.io/storage-class": "nfs",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "Job",
					Name:       Name,
					UID:        types.UID(Uid),
				},
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteMany,
			},
			StorageClassName: func() *string { s := "nfs"; return &s }(),
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
	return pvc
}
