package pprof

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

func Pprof() {
	runtime.GOMAXPROCS(1)              // 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪

	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}
