package pkg

func CreateResources(siteName string, deploymentName string, serviceName string, ingressName string, hostName string) {
	//deployment作成
	CreateDeployment(siteName, deploymentName)
	//service作成
	CreateService(siteName, serviceName)
	//ingress作成
	CreateIngress(ingressName, hostName, serviceName)
}

func DeleteResources(deploymentName string, serviceName string, ingressName string) {
	//deployment削除
	DeleteDeployment(deploymentName)
	//service削除
	DeleteService(serviceName)
	//ingress削除
	DeleteIngress(ingressName)
}
