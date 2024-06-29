package kubernetes

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deploymentを削除する処理
func DeleteDeployment(deploymentName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting deployment...")
	if err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete deployment: %v", err)
	}
	fmt.Println("Deleted deployment.")
	return nil
}

// serviceを削除する処理
func DeleteService(serviceName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting service...")
	if err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}
	fmt.Println("Deleted service.")
	return nil
}

// ingressを削除する処理
func DeleteIngress(ingressName string) error {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return err
	}

	ingressClient := clientset.NetworkingV1().Ingresses("default")
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting ingress...")
	if err := ingressClient.Delete(context.TODO(), ingressName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete ingress: %v", err)
	}
	fmt.Println("Deleted ingress.")
	return nil
}
