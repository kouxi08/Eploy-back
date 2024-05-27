package config

type KubeConfig struct {
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`
	IngressName    string `json:"ingressName"`
	HostName       string `json:"hostName"`
}

type Config struct {
	KubeConfig KubeConfig `json:"kubeConfig"`
}

