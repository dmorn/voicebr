package main

import (
	"context" "log"
	"os"
	"os/signal"

	"github.com/jecoz/callrelay"
	"github.com/jecoz/callrelay/enginevonage"
	"github.com/jecoz/callrelay/paths"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pidpath := paths.ServerPidPath()
	pid, _ := callrelay.LoadPid(pidpath)
	if pid > 0 {
		log.Printf("*** server (%d, %v) is already running", pid, pidpath)
		os.Exit(1)
	}
	if err := callrelay.SavePid(pidpath, os.Getpid()); err != nil {
		log.Printf("error * unable to save pid: %v", err)
		os.Exit(1)
	}
	defer os.Remove(pidpath)

	ev := enginevonage.New()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		log.Printf("*** signal %v received, canceling", <-interrupt)
		cancel()
	}()

	usr := make(chan os.Signal, 1)
	signal.Nofify(usr, syscall.SIGUSR1, SIGUSR2)
	go func() {
		for i := 0; true; i++ {
			log.Printf("*** signal %v received, reloading preferences (%d)", <-usr, i)
			var p *callrelay.Prefs
			if err := callrelay.LoadPrefs(paths.PregsPath(), &p); err != nil {
				log.Printf("error * reload preferences: %v", err)
				continue
			}
			if err := ev.SetPrefs(p); err != nil {
				log.Printf("error * reload preferences: %v", err)
			}
		}
	}()

	if err := ev.Run(ctx); err != nil {
		log.Printf("*** vonage broadcast engine: %v", err)
	}
}
