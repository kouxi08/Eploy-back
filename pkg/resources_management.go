package pkg

import (
	"fmt"
	"strconv"

	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/kouxi08/Eploy/utils"
)

type KanikoResult struct {
	HostName       string
	DeploymentName string
}

// kanikoを使ってbuild,pushをする際に使用するリソースをまとめたやつ
func CreateKanikoResouces(githubUrl string, appName string, targetPort string, envVars []kubernetes.EnvVar) (*KanikoResult, error) {
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

	//job作成
	jobName, jobUid, err := kubernetes.CreateJob(githubUrl, appName, registryName, envVars)
	if err != nil {
		return nil, err
	}
	//pvc作成
	if err := kubernetes.CreatePvc(jobName, jobUid, appName); err != nil {
		return nil, fmt.Errorf("failed to create PVC: %v", err)
	}
	errCh := make(chan error, 1)
	go func() {
		//jobの処理状況を監視
		errCh <- kubernetes.CheckJobCompletion(jobName)
	}()
	err = <-errCh
	if err != nil {
		return nil, err
	}

	//deployment作成
	err = kubernetes.CreateDeployment(appName, deploymentName, registryName, envVars)
	if err != nil {
		return nil, err
	}
	//service作成
	err = kubernetes.CreateService(appName, serviceName, targetPortInt)
	if err != nil {
		return nil, err
	}
	//ingress作成
	err = kubernetes.CreateIngress(ingressName, hostName, serviceName)
	if err != nil {
		return nil, err
	}

	result := &KanikoResult{
		HostName:       hostName,
		DeploymentName: deploymentName,
	}
	return result, nil
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
func GetStatusResources(deploymentName string) (status string, err error) {
	status, err = kubernetes.GetStatus(deploymentName)
	return status, err
}
