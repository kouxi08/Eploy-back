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
	//deploymentの定義
	deployment := DeploymentDefinition(app, deploymentName)

	//k8sに送信
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
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
	//serviceの定義
	service := ServiceDefinition(app, serviceName)

	//k8sに送信
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
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
	//ingressの定義
	ingress := IngressDefinition(ingressName, hostName, serviceName)

	//k8sに送信
	ingressClient := clientset.NetworkingV1().Ingresses("default")
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ingress: %v", err)
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// kanikoのjobを生成する処理
func CreateJob() (string, string, error) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return "", "", err
	}
	//jobの定義
	job := JobDefinition()

	//k8sに送信
	jobClient := clientset.BatchV1().Jobs("default")
	result, err := jobClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to create job: %v", err)
	}
	uid := string(result.UID)
	name := result.Name
	fmt.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())
	return name, uid, nil
}

// pvcを作成する処理
func CreatePvc(Name string, Uid string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}
	//pvcの定義
	pvc := PvcDefinition(Name, Uid)

	// PVCを作成
	pvcClient := clientset.CoreV1().PersistentVolumeClaims("default")
	result, err := pvcClient.Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created PVC %q.\n", result.GetObjectMeta().GetName())
	return err
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
