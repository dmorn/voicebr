package nexmo_test

import (
	"context"
	"time"
	"testing"

	"github.com/jecoz/voicebr/nexmo"
)

func TestLimit_exceed(t *testing.T) {
	// Limiter that allows 3 req/sec
	l := nexmo.NewLimiter(3)

	// Consume the bucket
	ctx := context.TODO()
	for i := 0; i < 3; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("%d: Unexpected limiter error: %v", i, err)
		}
	}

	// Now it should fail
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err == nil {
		t.Fatal("Limiter did not block")
	}
}

func TestLimit(t *testing.T) {
	l := nexmo.NewLimiter(1)

	ctx := context.TODO()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Unexpected limiter error: %v", err)
	}
	<-time.After(time.Second)
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Unexpected limiter error: %v", err)
	}
}
