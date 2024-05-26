package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// deploymentを作成する処理
func CreateDeployment(app string, deploymentName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

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
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// serviceを作成する処理
func CreateService(app string, serviceName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

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
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service%q.\n", result.GetObjectMeta().GetName())
	return nil
}

// ingressを作成する処理
func CreateIngress(ingressName string, hostName string, serviceName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	ingressClient := clientset.NetworkingV1().Ingresses("default")

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
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// kanikoのpodを生成
func CreateKaniko() error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	podClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko",
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
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
			Containers: []apiv1.Container{
				{
					Name:  "test-kaniko",
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
		},
	}
	result, err := podClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// pod内のlogを取得
func GetPodLog(podName string) (string, error) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}
	namespace := "default" // Specify the namespace
	fmt.Println(podName)

	podLogOpts := apiv1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		log.Println("GetPodLog:")
		log.Println(err)
		return "", err
	}

	defer podLogs.Close()
	var sb strings.Builder
	buf := make([]byte, 2000)
	for {
		numBytes, err := podLogs.Read(buf)
		if numBytes == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		sb.Write(buf[:numBytes])
	}
	logOutput := sb.String()
	fmt.Println(logOutput)
	return logOutput, nil
}
