package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"unsafe"
)

var K8s []K8sConf

func initK8sCreds() error {
	k8sBytes := unsafe.Slice(unsafe.StringData(Env.K8s), len(Env.K8s))
	err := yaml.Unmarshal(k8sBytes, &K8s)
	if err != nil {
		return err
	}

	if len(K8s) == 0 {
		return errors.New("invalid empty K8s Creds")
	}
	return nil
}

type K8sConf struct {
	Backend string `yaml:"backend"`
	Token   string `yaml:"token"`
}
