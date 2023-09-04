package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"unsafe"
)

// Workloads namespace ==> workload ==> conf
var Workloads []map[string]map[string]WorkloadConf

func initWorkloads() error {
	workloadBytes := unsafe.Slice(unsafe.StringData(Env.Workloads), len(Env.Workloads))
	err := yaml.Unmarshal(workloadBytes, &Workloads)
	if err != nil {
		return err
	}

	if len(Workloads) == 0 {
		return errors.New("invalid empty workloads")
	}
	return nil
}

type WorkloadConf struct {
	Image string `yaml:"image"`
	Wait  bool   `yaml:"wait"`
}
