package utils

import (
	"os"
	"os/signal"
	"syscall"
)

/*
 @Author: 71made
 @Date: 2023/02/14 23:31
 @ProductName: signal.go
 @Description:
*/

func DealSignal(do func()) {
	// 关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		do()
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
}
