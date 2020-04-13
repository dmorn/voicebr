package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jecoz/voicebr"
	"github.com/jecoz/voicebr/enginevonage"
	"github.com/jecoz/voicebr/paths"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pidpath := paths.ServerPidPath()
	pid, _ := voicebr.LoadPid(pidpath)
	if pid > 0 {
		log.Printf("*** server (%d, %v) is already running", pid, pidpath)
		os.Exit(1)
	}
	if err := voicebr.SavePid(pidpath, os.Getpid()); err != nil {
		log.Printf("error * unable to save pid: %v", err)
		os.Exit(1)
	}
	defer os.Remove(pidpath)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		log.Printf("*** signal %v received, canceling", <-sig)
		cancel()
	}()

	if err := enginevonage.New().Run(ctx); err != nil {
		log.Printf("*** vonage broadcast engine: %v", err)
	}
}
