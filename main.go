package main

import (
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/shivamhw/piro-go/app"
)

func main(){
	runBot()
}

func runBot() {
	for {
		func() {
			defer func(){
				if r := recover(); r!=nil {
					slog.Error("panic occured", "err", r)
					debug.Stack()
					slog.Info("recovering in 10s......")
					time.Sleep(10 * time.Second)
				}
			}()
			app.Run()
		}()
	}
}
