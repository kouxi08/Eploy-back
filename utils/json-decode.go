package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	KubeManifest KubeManifest `json:"KubeManifest"`
	DNSRecords   DNSRecord    `json:"dns_records"`
}

type KubeManifest struct {
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`
	IngressName    string `json:"ingressName"`
	HostName       string `json:"hostName"`
	RegistryName   string `json:"registryName"`
	TargetPort     string `json:"targetPort"`
}

type DNSRecord struct {
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}
	return &config, nil
}
