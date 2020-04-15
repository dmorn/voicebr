package enginevonage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jecoz/callrelay"
	"github.com/jecoz/callrelay/vonage"
)

type httpError struct {
	Code int
	Err  error
}

func (e *httpError) Error() string {
	return e.Err.Error()
}

func Error(w http.ResponseWriter, err error) {
	log.Printf("error * %v", err)

	code := http.StatusInternalServerError
	var herr *httpError
	if errors.As(err, &herr) {
		code = herr.Code
	}
	http.Error(w, err.Error(), code)
}

// ReturnHandlerFunc describes an HandlerFunc which returns the errors
// it encounters. It is a http.Handler implementation. Errors that are
// of type httpError should be casted to unwrap error's status code,
// 500 should be used for the other cases.
type ReturnHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Implementation of http.Handler.
func (h ReturnHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		Error(w, err)
		return
	}
}

type Engine struct {
	Prefs  *callrelay.Prefs
	Config *vonage.Config
}

func (e *Engine) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", e.Prefs.Port)
	callback := fmt.Sprintf("%v:%d/voice/broadcast/event", e.Prefs.ExternalOrigin, e.Prefs.Port)
	curl, err := url.Parse(callback)
	if err != nil {
		return fmt.Errorf("engine run: %w", err)
	}

	recMux := vonage.NewVoiceWebhookMux(&vonage.VoiceWebhook{
		Answer: NewRecordHandlerFunc(curl.String(), e.Prefs, e.Config),
		Event:  ReturnHandlerFunc(Event),
	})
	brMux := vonage.NewVoiceWebhookMux(vonage.NewWebook())

	mux := http.NewServeMux()
	mux.Handle("/voice/record/", http.StripPrefix("/voice/record", recMux))
	mux.Handle("/voice/broadcast/", http.StripPrefix("/voice/broadcast", brMux))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	done := make(chan error)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		done <- srv.Shutdown(ctx)
	}()

	log.Printf("server listening on addr: %v", addr)
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("engine run: %w", err)
	}
	if err := <-done; err != nil {
		return fmt.Errorf("engine shutdown: %w", err)
	}
	return nil
}
