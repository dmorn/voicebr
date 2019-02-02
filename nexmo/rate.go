package nexmo

import (
	"context"

	"golang.org/x/time/rate"
)

type Limiter struct {
	internal *rate.Limiter
}

// NewLimiter creates a new Limiter intance that
// allows `r` events per second to happen.
// It is a wrapper around rate.Limiter, configured
// with an burst of `r` and Limit of `r`.
func NewLimiter(r int) *Limiter {
	return &Limiter{
		internal: rate.NewLimiter(rate.Limit(r), r),
	}
}

// Wait blocks until the caller is allowed to perform
// a request acoording to the limiter's configuration.
func (l *Limiter) Wait(ctx context.Context) error {
	return l.internal.Wait(ctx)
}
