package config

type KubeConfig struct {
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`
	IngressName    string `json:"ingressName"`
	HostName       string `json:"hostName"`
	RegistryName   string `json:"registryName"`
	TargetPort     string `json:"targetPort"`
}

type Config struct {
	KubeConfig KubeConfig `json:"kubeConfig"`
}
