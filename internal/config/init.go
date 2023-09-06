package config

import (
	"github.com/caarlos0/env/v6"
	"log/slog"
	"os"
)

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
