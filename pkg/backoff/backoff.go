package backoff

import (
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"log/slog"
	"time"
)

// Backoff 错误重试 积分退避算法
type Backoff struct {
	c Conf

	retryDelay time.Duration
}

type Conf struct {
	Name    string
	Content func() error
	// 最大重试等待时间
	MaxRetryDelay time.Duration
	// 最大重试次数，0 为不设限
	MaxRetry uint8
}

func New(c Conf) Backoff {
	if c.Name == "" {
		c.Name = "UNKNOWN"
	}
	if c.Content == nil {
		panic("content function required")
	}
	if c.MaxRetryDelay == 0 {
		c.MaxRetryDelay = time.Minute * 20
	}

	return Backoff{
		c:          c,
		retryDelay: time.Second * 2,
	}
}

func (a Backoff) Start() {
	go a.Worker()
}

// Worker
// 请注意,此处使用的是普通接收器,当 worker 重新运行时参数会被重置
func (a Backoff) Worker() {
	quitAfterCount := a.c.MaxRetry != 0

	for {
		errChan := make(chan error)
		go func() {
			defer func() {
				if p := tool.Recover(); p != nil {
					errChan <- errors.New(fmt.Sprint(p))
				}
			}()
			errChan <- a.c.Content()
		}()
		if err := <-errChan; err == nil {
			break
		}

		if quitAfterCount {
			if a.c.MaxRetry == 0 {
				slog.Error("backoff deadline reached",
					slog.String("name", a.c.Name),
				)
				break
			}
			a.c.MaxRetry--
		}

		slog.Info("backoff retry...",
			slog.String("name", a.c.Name),
			slog.Duration("wait", a.retryDelay),
		)

		time.Sleep(a.retryDelay)

		if a.retryDelay < a.c.MaxRetryDelay {
			a.retryDelay = a.retryDelay << 1
			if a.retryDelay > a.c.MaxRetryDelay {
				a.retryDelay = a.c.MaxRetryDelay
			}
		}
	}
}
