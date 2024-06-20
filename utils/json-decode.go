package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/kouxi08/Eploy/pkg/kubernetes"
)

func LoadConfig(filename string) (*kubernetes.KubeConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config kubernetes.KubeConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}
	return &config, nil
}
