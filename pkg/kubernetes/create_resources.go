package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strings"

	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deploymentを作成する処理
func CreateDeployment(app string, deploymentName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}
	deployment := Deployment_definition(app, deploymentName)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
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
	service := Service_definition(app, serviceName)

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
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
	ingress := Ingress_definition(ingressName, hostName, serviceName)

	ingressClient := clientset.NetworkingV1().Ingresses("default")
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// kanikoのjobを生成
func CreateKaniko() error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}
	job := Job__definition()

	jobClient := clientset.BatchV1().Jobs("default")
	result, err := jobClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// pod内のlogを取得する処理
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
