package pkg

import (
	"fmt"
	"strconv"

	"github.com/kouxi08/Eploy/config"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
)

// アプリケーションを作成する際に動作させるリソースをまためたやつ
func CreateResources(siteName string, targetPort string) {
	config, _ := config.LoadConfig("config.json")

	// deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)
	hostName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.HostName)
	// registryName := fmt.Sprintf("%s%s", config.KubeConfig.RegistryName, siteName)
	targetPortInt, _ := strconv.Atoi(targetPort)

	//deployment作成
	// kubernetes.CreateDeployment(siteName, deploymentName, registryName, )
	//service作成
	kubernetes.CreateService(siteName, serviceName, targetPortInt)
	//ingress作成
	kubernetes.CreateIngress(ingressName, hostName, serviceName)
}

// kanikoを使ってbuild,pushをする際に使用するリソースをまとめたやつ
func CreateKanikoResouces(githubUrl string, appName string, targetPort string, envVars []kubernetes.EnvVar) error {
	config, _ := config.LoadConfig("config.json")

	deploymentName := fmt.Sprintf("%s%s", appName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", appName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", appName, config.KubeConfig.IngressName)
	hostName := fmt.Sprintf("%s%s", appName, config.KubeConfig.HostName)
	registryName := fmt.Sprintf("%s%s", config.KubeConfig.RegistryName, appName)
	targetPortInt, err := strconv.Atoi(targetPort)
	if err != nil {
		return err
	}

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
	kubernetes.CreateDeployment(appName, deploymentName, registryName, envVars)
	//service作成
	kubernetes.CreateService(appName, serviceName, targetPortInt)
	//ingress作成
	kubernetes.CreateIngress(ingressName, hostName, serviceName)

	return nil
}

// アプリケーションを削除する際に動作させるリソースを定義したやつ
func DeleteResources(siteName string) {
	config, _ := config.LoadConfig("config.json")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)

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

// podのステータスを確認するやつ
func GetStatusResources(deploymentName string) (status string,err error){
	status,err = kubernetes.GetStatus(deploymentName)
	return status,err
}
