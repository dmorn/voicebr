package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jecoz/voicebr/vonage"
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

type ReturnHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h ReturnHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		Error(w, err)
		return
	}
}

func Answer(w http.ResponseWriter, r *http.Request) error {
	p, err := vonage.ParseVoiceAnswerRequest(r)
	if err != nil {
		return &httpError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("answer request: %w", err),
		}
	}

	log.Printf("*** reveived call (%v) from: %v", p.UUID, p.From)

	nccos := []interface{}{vonage.NewTalkControl("Hi from Jecoz")}
	if err = json.NewEncoder(w).Encode(nccos); err != nil {
		return fmt.Errorf("answer response: %w", err)
	}
	return nil
}

func Event(w http.ResponseWriter, r *http.Request) error {
	p, err := vonage.ParseVoiceEventRequest(r)
	if err != nil {
		return &httpError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("answer request: %w", err),
		}
	}

	log.Printf("event: %+v", p)
	return nil
}

func main() {
	vw := &vonage.VoiceWebhook{
		Answer: ReturnHandlerFunc(Answer),
		Event:  ReturnHandlerFunc(Event),
	}
	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: vonage.NewVoiceWebhookMux(vw),
	}
	log.Printf("server listening on addr: %v", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("*** listener: %v", err)
	}
}
