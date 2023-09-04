package config

import (
	"github.com/caarlos0/env/v6"
	"log/slog"
	"os"
)

var Env EnvConf

func init() {
	err := env.Parse(&Env, env.Options{
		Prefix: "INPUT_",
	})
	if err != nil {
		slog.Error("read env failed", slog.Any("err", err))
		os.Exit(1)
	}

	if err = initK8sCreds(); err != nil {
		slog.Error("decode K8s Creds failed", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("K8s Creds loaded", slog.Int("num", len(K8s)))

	if err = initWorkloads(); err != nil {
		slog.Error("decode Workloads failed", slog.Any("err", err))
		os.Exit(1)
	}
}

type EnvConf struct {
	K8s       string `env:"K8S,required"`       // 集群 api 凭据
	Workloads string `env:"WORKLOADS,required"` // 目标工作负载
}
