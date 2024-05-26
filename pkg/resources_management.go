package pkg

import "github.com/kouxi08/Eploy/pkg/kubernetes"

func CreateResources(siteName string, deploymentName string, serviceName string, ingressName string, hostName string) {
	//deployment作成
	kubernetes.CreateDeployment(siteName, deploymentName)
	//service作成
	kubernetes.CreateService(siteName, serviceName)
	//ingress作成
	kubernetes.CreateIngress(ingressName, hostName, serviceName)
}

func CreateKanikoResouces() {
	kubernetes.CreateJob()
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
