package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jecoz/voley"
	"github.com/jecoz/voley/enginevonage"
	"github.com/jecoz/voley/vonage"
)

func logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func errorf(format string, args ...interface{}) {
	logf("error * "+format, args...)
}

// deferred functions will not be called if this function is.
func exitf(s int, format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(s)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p *voley.Prefs
	if err := p.Load(paths.Prefs()); err != nil {
		exitf(1, err.Error())
	}

	// Instantiate the Vonage engine.
	ev := enginevonage.Engine{Prefs: &p}

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
