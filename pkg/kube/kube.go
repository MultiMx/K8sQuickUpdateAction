package kube

import "strings"

type Kube struct {
	Conf *Config
}

func New(conf *Config) *Kube {
	conf.Backend = strings.TrimSuffix(conf.Backend, "/")
	if conf.Cluster == "" {
		conf.Cluster = "local"
	}
	return &Kube{
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
