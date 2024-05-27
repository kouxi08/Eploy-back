package pkg

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	// "errors"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset, nil
}

func GetKubernetesNodes() {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for y, nodes := range nodes.Items {
		fmt.Printf("[%d] %s\n", y, nodes.GetName())
	}

}

// deploymentを作成する処理
func CreateDeployment(app string, deploymentName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
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
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

// serviceを作成する処理
func CreateService(app string, serviceName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
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

	// Create Service
	fmt.Println("Creating service...")
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service%q.\n", result.GetObjectMeta().GetName())
}

// ingressを作成する処理
func CreateIngress(ingressName string, hostName string, serviceName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
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
	fmt.Println("Creating ingress...")
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
}

// deploymentを削除する処理
func DeleteDeployment(deploymentName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting deployment...")
	if err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

// serviceを削除する処理
func DeleteService(serviceName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting service...")
	if err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service.")
}

// ingressを削除する処理
func DeleteIngress(ingressName string) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	ingressClient := clientset.NetworkingV1().Ingresses("default")
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting ingress...")
	if err := ingressClient.Delete(context.TODO(), ingressName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted ingress.")
}

func GetPodLog(podName string) (string ,error){
    clientset, err := NewKubernetesClient()
    if err != nil {
        log.Fatal(err)
    }
    namespace := "default"      // Specify the namespace
	fmt.Println(podName)

    podLogOpts := apiv1.PodLogOptions{}
    req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
    podLogs, err := req.Stream(context.TODO())
    if err != nil {
		log.Println("GetPodLog:")
        log.Println(err)
		return "",err
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
    return logOutput,nil
}
func int32Ptr(i int32) *int32 { return &i }
