package controllers

import (
	"context"
	"errors"
	"github.com/MultiMx/K8sQuickUpdateAction/internal/config"
	"github.com/MultiMx/K8sQuickUpdateAction/pkg/backoff"
	"github.com/MultiMx/K8sQuickUpdateAction/pkg/kube"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type DeployWorkChain struct {
	Prev *DeployWorkChain
	F    func()
}

func UpdateWorkloads() {
	var apis = make([]*kube.Kube, len(config.K8s))
	i := 0
	for name, conf := range config.K8s {
		apis[i] = kube.New(&conf)
		apis[i].Name = name
		i++
	}

	var errCount atomic.Uint32

	for i, step := range config.Workloads {
		slog.Info("---------", "index", i)

		var wg sync.WaitGroup
		var chain *DeployWorkChain
		for namespace, workloads := range step {
			for workload, conf := range workloads {
				conf := conf
				for _, kubeApi := range apis {
					kubeApi := kubeApi
					wg.Add(1)
					operator := kubeApi.NewWorkload(namespace, workload)
					var deployWork = DeployWorkChain{
						Prev: chain,
						F: func() {
							defer wg.Done()
							logger := slog.With(
								slog.String("kube", kubeApi.Name),
								slog.String("namespace", operator.Namespace),
								slog.String("workload", operator.Workload),
							)

							err := backoff.New(backoff.Conf{
								Logger: logger.With("action", "set image"),
								Content: func() error {
									err := operator.SetImage(conf.Image)
									if err != nil {
										errCount.Add(1)
										logger.Warn("set image failed", "err", err)
									}
									return err
								},
								MaxRetryDelay: time.Second * 5,
								MaxRetry:      8,
							}).Run()
							if err != nil {
								logger.Error("set image backoff failed", "err", err)
								return
							}

							if conf.Wait {
								err = WaitDeploymentAvailable(logger, operator)
								if err != nil {
									errCount.Add(1)
									logger.Error("wait for full available failed", "err", err)
									return
								}
							}

							logger.Info("success")
						},
					}
					chain = &deployWork
				}
			}
		}

		for chain != nil {
			go chain.F()
			chain = chain.Prev
		}
		wg.Wait()
	}

	if errCount.Load() != 0 {
		slog.Error("completed with error")
		os.Exit(1)
	} else {
		slog.Info("completed successfully")
	}
}

// WaitDeploymentAvailable 5 times error or 5 minutes
func WaitDeploymentAvailable(logger *slog.Logger, workload kube.Workload) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	var count uint8
	for {
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			return errors.New("wait for full available timeout")
		default:
			ok, err := workload.FullAvailable()
			if err != nil {
				logger.Warn("check available status failed", "err", err)
				count++
				if count >= 5 {
					return errors.New("maximum number of error retries reached")
				}
			} else if ok {
				return nil
			}
		}
	}
}
