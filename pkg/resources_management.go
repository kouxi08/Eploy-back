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

func CreateResourceNames(appName string, targetPort string) (*ResponseNames, error) {
	config, _ := utils.LoadConfig("config.json")
	deploymentName := fmt.Sprintf("%s%s", appName, config.KubeManifest.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", appName, config.KubeManifest.ServiceName)
	ingressName := fmt.Sprintf("%s%s", appName, config.KubeManifest.IngressName)
	hostName := fmt.Sprintf("%s%s", appName, config.KubeManifest.HostName)
	registryName := fmt.Sprintf("%s%s", config.KubeManifest.RegistryName, appName)
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
func CreateResouces(githubUrl string, appName string, resourceNames *ResponseNames, envVars []kubernetes.EnvVar) error {
	registryName := resourceNames.RegistryName
	deploymentName := resourceNames.DeploymentName
	serviceName := resourceNames.ServiceName
	ingressName := resourceNames.IngressName
	hostName := resourceNames.HostName
	targetPortInt := resourceNames.TargetPortInt

	//job作成
	jobName, jobUid, err := kubernetes.CreateJob(githubUrl, appName, registryName, envVars)
	if err != nil {
		return err
	}
	//pvc作成
	if err := kubernetes.CreatePvc(jobName, jobUid, appName); err != nil {
		return fmt.Errorf("failed to create PVC: %v", err)
	}
	errCh := make(chan error, 1)
	go func() {
		//jobの処理状況を監視
		errCh <- kubernetes.CheckJobCompletion(jobName)
	}()
	err = <-errCh
	if err != nil {
		return err
	}

	//deployment作成
	err = kubernetes.CreateDeployment(appName, deploymentName, registryName, envVars)
	if err != nil {
		return err
	}
	//service作成
	err = kubernetes.CreateService(appName, serviceName, targetPortInt)
	if err != nil {
		return err
	}
	//ingress作成
	err = kubernetes.CreateIngress(ingressName, hostName, serviceName)
	if err != nil {
		return err
	}
	return nil
}

// アプリケーションを削除する際に動作させるリソースを定義したやつ
func DeleteResources(siteName string) error {
	utils, _ := utils.LoadConfig("config.json")

	deploymentName := fmt.Sprintf("%s%s", siteName, utils.KubeManifest.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, utils.KubeManifest.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, utils.KubeManifest.IngressName)

	//deployment削除
	err := kubernetes.DeleteDeployment(deploymentName)
	if err != nil {
		return err
	}
	//service削除
	err = kubernetes.DeleteService(serviceName)
	if err != nil {
		return err
	}
	//ingress削除
	err = kubernetes.DeleteIngress(ingressName)
	if err != nil {
		return err
	}

	return nil
}

func GetLogPodResources(podName string) (message string, err error) {
	message, err = kubernetes.GetPodLog(podName)
	return
}

// podのステータスを確認するやつ
func GetStatusResources(deploymentName string, jobsName string) (status string, err error) {
	status, err = kubernetes.GetDeploymentStatus(deploymentName)
	if err != nil {
		status, err = kubernetes.GetJobsStatus(jobsName)
	}
	if err != nil {
		return "", err
	}
	return status, err
}
