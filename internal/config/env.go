package config

var Env EnvConf

type EnvConf struct {
	K8s       string `env:"K8S,required"`       // 集群 api 凭据
	Workloads string `env:"WORKLOADS,required"` // 目标工作负载
}
