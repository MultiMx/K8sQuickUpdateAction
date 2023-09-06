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
	Logger *slog.Logger

	Content func() error
	// 最大重试等待时间
	MaxRetryDelay time.Duration
	// 最大重试次数，0 为不设限
	MaxRetry uint8
}

func New(c Conf) Backoff {
	if c.MaxRetryDelay == 0 {
		c.MaxRetryDelay = time.Minute * 20
	}

	return Backoff{
		c:          c,
		retryDelay: time.Second,
	}
}

// Run
// 请注意,此处使用的是普通接收器,当 worker 重新运行时参数会被重置
func (a Backoff) Run() error {
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
				return errors.New("backoff deadline reached")
			}
			a.c.MaxRetry--
		}

		a.c.Logger.Info("backoff retry...",
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

	return nil
}
