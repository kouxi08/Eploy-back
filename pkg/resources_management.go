package pkg

import (
	"fmt"

	"github.com/kouxi08/Eploy/pkg/kubernetes"
)

// アプリケーションを作成する際に動作させるリソースをまためたやつ
func CreateResources(siteName string, deploymentName string, serviceName string, ingressName string, hostName string) {
	//deployment作成
	kubernetes.CreateDeployment(siteName, deploymentName)
	//service作成
	kubernetes.CreateService(siteName, serviceName)
	//ingress作成
	kubernetes.CreateIngress(ingressName, hostName, serviceName)
}

// kanikoを使ってbuild,pushをする際に使用するリソースをまとめたやつ
func CreateKanikoResouces(githubUrl string, appName string, envVars []kubernetes.EnvVar) error {

	//job作成
	jobName, jobUid, err := kubernetes.CreateJob(githubUrl, appName, envVars)
	if err != nil {
		return err
	}
	//jobの処理状況を監視
	go kubernetes.CheckJobCompletion(jobName)

	//pvc作成
	if err := kubernetes.CreatePvc(jobName, jobUid, appName); err != nil {
		return fmt.Errorf("failed to create PVC: %v", err)
	}

	return nil
}

// アプリケーションを削除する際に動作させるリソースを定義したやつ
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
