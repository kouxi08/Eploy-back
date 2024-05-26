package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Deployment_definition(app string, deploymentName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "kouxi00/portfolio",
						},
					},
				},
			},
		},
	}
	return deployment
}

func Service_definition(app string, serviceName string) *apiv1.Service {

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(3000),
				},
			},
			Selector: map[string]string{
				"app": app,
			},
		},
	}
	return service
}

func Ingress_definition(ingressName string, hostName string, serviceName string) *networkingv1.Ingress {

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

func Job__definition() *batchv1.Job {
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
					Containers: []apiv1.Container{
						{
							Name:  "kaniko",
							Image: "gcr.io/kaniko-project/executor:latest",
							Args: []string{
								"--dockerfile=/workspace/Dockerfile",
								"--context=dir:///workspace",
								"--no-push",
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "dockerfile-storage",
									MountPath: "/workspace",
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
									ClaimName: "local-claim",
								},
							},
						},
					},
					InitContainers: []apiv1.Container{
						{
							Name:  "init-gitclone",
							Image: "alpine/git:2.43.4",
							Command: []string{
								"git",
								"clone",
								"https://github.com/kouxi08/Eploy-back.git",
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
				},
			},
			BackoffLimit:            int32Ptr(0),
			TTLSecondsAfterFinished: int32Ptr(30),
		},
	}
	return job
}
