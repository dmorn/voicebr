package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jecoz/voicebr/enginevonage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
