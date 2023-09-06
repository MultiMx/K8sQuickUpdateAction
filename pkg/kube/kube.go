package kube

import "strings"

type Kube struct {
	Name string
	Conf *Config
}

func New(name string, conf *Config) Kube {
	conf.Backend = strings.TrimSuffix(conf.Backend, "/")
	if conf.Cluster == "" {
		conf.Cluster = "local"
	}
	return Kube{
		Name: name,
		Conf: conf,
	}
}

func (a Kube) NewWorkload(namespace, workload string) Workload {
	return Workload{
		kube:      a,
		Namespace: namespace,
		Workload:  workload,
	}
}
