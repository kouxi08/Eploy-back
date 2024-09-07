package pkg

import (
	"fmt"
	"strconv"

	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/kouxi08/Eploy/utils"
)

type ResponseNames struct {
	HostName       string
	DeploymentName string
	ServiceName    string
	IngressName    string
	RegistryName   string
	TargetPortInt  int
}

func CreateResourceNames(kube *utils.KubeManifest, appName string, targetPort string) (*ResponseNames, error) {
	deploymentName := fmt.Sprintf("%s%s", appName, kube.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", appName, kube.ServiceName)
	ingressName := fmt.Sprintf("%s%s", appName, kube.IngressName)
	hostName := fmt.Sprintf("%s%s", appName, kube.HostName)
	registryName := fmt.Sprintf("%s%s", kube.RegistryName, appName)
	targetPortInt, err := strconv.Atoi(targetPort)
	if err != nil {
		return nil, err
	}

	result := &ResponseNames{
		HostName:       hostName,
		DeploymentName: deploymentName,
		ServiceName:    serviceName,
		IngressName:    ingressName,
		RegistryName:   registryName,
		TargetPortInt:  targetPortInt,
	}

	return result, nil
}

// kanikoを使ってbuild,pushをする際に使用するリソースをまとめたやつ
func CreateResouces(k *kubernetes.KubernetesApp, githubUrl string, appName string, resourceNames *ResponseNames, envVars []kubernetes.EnvVar) error {
	registryName := resourceNames.RegistryName
	deploymentName := resourceNames.DeploymentName
	serviceName := resourceNames.ServiceName
	ingressName := resourceNames.IngressName
	hostName := resourceNames.HostName
	targetPortInt := resourceNames.TargetPortInt

	//job作成
	jobName, jobUid, err := k.CreateJob(githubUrl, appName, registryName, envVars)
	if err != nil {
		return err
	}
	//pvc作成
	if err := k.CreatePvc(jobName, jobUid, appName); err != nil {
		return fmt.Errorf("failed to create PVC: %v", err)
	}
	errCh := make(chan error, 1)
	go func() {
		//jobの処理状況を監視
		errCh <- k.CheckJobCompletion(jobName)
	}()
	err = <-errCh
	if err != nil {
		return err
	}

	//deployment作成
	err = k.CreateDeployment(appName, deploymentName, registryName, envVars)
	if err != nil {
		return err
	}
	//service作成
	err = k.CreateService(appName, serviceName, targetPortInt)
	if err != nil {
		return err
	}
	//ingress作成
	err = k.CreateIngress(ingressName, hostName, serviceName)
	if err != nil {
		return err
	}
	return nil
}

// アプリケーションを削除する際に動作させるリソースを定義したやつ
func DeleteResources(k *kubernetes.KubernetesApp, kube *utils.KubeManifest, siteName string) error {

	deploymentName := fmt.Sprintf("%s%s", siteName, kube.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, kube.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, kube.IngressName)

	//deployment削除
	err := k.DeleteDeployment(deploymentName)
	if err != nil {
		return err
	}
	//service削除
	err = k.DeleteService(serviceName)
	if err != nil {
		return err
	}
	//ingress削除
	err = k.DeleteIngress(ingressName)
	if err != nil {
		return err
	}

	return nil
}

func GetLogPodResources(k *kubernetes.KubernetesApp, podName string) (message string, err error) {
	message, err = k.GetPodLog(podName)
	return
}

// podのステータスを確認するやつ
func GetStatusResources(k *kubernetes.KubernetesApp, deploymentName string, jobsName string) (status string, err error) {
	status, err = k.GetDeploymentStatus(deploymentName)
	if err != nil {
		status, err = k.GetJobsStatus(jobsName)
	}
	if err != nil {
		return "", err
	}
	return status, err
}
