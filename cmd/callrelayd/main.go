package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jecoz/callrelay"
	"github.com/jecoz/callrelay/vonage"
	"github.com/jecoz/callrelay/enginevonage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prefspath := "callrelayprefs.hujson"
	configpath := "vonageconf.hujson"

	var p callrelay.Prefs
	if err := callrelay.LoadPrefs(prefspath, &p); err != nil {
		log.Printf("error * %v", err)
		os.Exit(1)
	}
	var c vonage.Config
	if err := vonage.LoadConfig(configpath, &c); err != nil {
		log.Printf("error * %v", err)
		os.Exit(1)
	}

	// Instantiate the Vonage engine.
	ev := enginevonage.Engine{
		Prefs: &p,
		Config: &c,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		log.Printf("*** signal %v received, canceling", <-interrupt)
		cancel()
	}()

	if err := ev.Run(ctx); err != nil {
		log.Printf("*** vonage broadcast engine: %v", err)
	}
}
