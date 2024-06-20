package kubernetes

type KubeManifest struct {
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`
	IngressName    string `json:"ingressName"`
	HostName       string `json:"hostName"`
	RegistryName   string `json:"registryName"`
	TargetPort     string `json:"targetPort"`
}

type KubeConfig struct {
	KubeManifest KubeManifest `json:"KubeManifest"`
}
