package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunTimer(num int, uint Unit, fn func()) {
	go func() {
		sigChan := make(chan os.Signal, 1)
		// run once right now
		// 捕获指定的系统信号
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		t := time.NewTimer(time.Second * 1)
		fn()
		for {
			now := time.Now()
			next := now.Add(ParseDuration(num, uint))
			next = TruncateTime(next, uint)
			t.Reset(next.Sub(now))
			select {
			case <-t.C:
				fn()
			case <-sigChan:
				fmt.Println("timer is exist...")
				return
			}
		}
	}()
}
