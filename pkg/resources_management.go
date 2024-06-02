package pkg

import (
	"fmt"

	"github.com/kouxi08/Eploy/pkg/kubernetes"
)

func CreateResources(siteName string, deploymentName string, serviceName string, ingressName string, hostName string) {
	//deployment作成
	kubernetes.CreateDeployment(siteName, deploymentName)
	//service作成
	kubernetes.CreateService(siteName, serviceName)
	//ingress作成
	kubernetes.CreateIngress(ingressName, hostName, serviceName)
}

func CreateKanikoResouces(githubUrl string, appName string, envVars []kubernetes.EnvVar) error {

	//pvc作成
	jobName, jobUid, err := kubernetes.CreateJob(githubUrl, appName, envVars)
	if err != nil {
		return err
	}
	go kubernetes.CheckJobCompletion(jobName)

	//job作成
	if err := kubernetes.CreatePvc(jobName, jobUid, appName); err != nil {
		return fmt.Errorf("failed to create PVC: %v", err)
	}

	return nil
}

func DeleteResources(deploymentName string, serviceName string, ingressName string) {
	//deployment削除
	kubernetes.DeleteDeployment(deploymentName)
	//service削除
	kubernetes.DeleteService(serviceName)
	//ingress削除
	kubernetes.DeleteIngress(ingressName)
}

func GetLogPodResources(podName string) (message string, err error) {
	message, err = kubernetes.GetPodLog(podName)
	return
}
