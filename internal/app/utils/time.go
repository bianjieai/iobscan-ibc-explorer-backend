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
		signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 创建一个timer
		t := time.NewTimer(time.Second * 5)
		fn()
		for {
			now := time.Now()
			next := now.Add(ParseDuration(num, uint))
			next = TruncateTime(next, uint)
			//使用 time.Reset 重置 timer，重复利用 timer
			t.Reset(next.Sub(now))
			select {
			case <-t.C:
				fn()
			case <-sigChan:
				t.Stop()
				fmt.Println("timer is exist...")
				return
			}
		}
	}()
}
