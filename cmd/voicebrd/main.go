package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jecoz/voicebr/vonage"
)

type codeError struct {
	Code int
	Err  error
}

func (e *codeError) Error() string {
	return e.Err.Error()
}

func Error(w http.ResponseWriter, err error) {
	log.Printf("error * %v", err)

	code := http.StatusInternalServerError
	var cerr *codeError
	if errors.As(err, &cerr) {
		code = cerr.Code
	}
	http.Error(w, err.Error(), code)
}

func returnAnswer(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read answer request: %v", err)
	}
	var p vonage.VoiceAnswerRequest
	if err = json.Unmarshal(b, &p); err != nil {
		return &codeError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("unable to unmarshal answer request: %w", err),
		}
	}

	log.Printf("*** reveived call (%v) from: %v", p.UUID, p.From)

	ncco := vonage.NewTalkControl("Hi from Jecoz")
	if err = json.NewEncoder(w).Encode(&ncco); err != nil {
		return fmt.Errorf("unable to send response: %w", err)
	}
	return nil
}

func Answer(w http.ResponseWriter, r *http.Request) {
	if err := returnAnswer(w, r); err != nil {
		Error(w, err)
	}
	// Just need to write a response, http.StatusOK will be
	// send automatically.
}

func returnEvent(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read event request: %v", err)
	}
	var p vonage.VoiceEventRequest
	if err = json.Unmarshal(b, &p); err != nil {
		return &codeError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("unable to unmarshal event request: %w", err),
		}
	}

	log.Printf("event: %v", p)
	w.WriteHeader(http.StatusOK)
	return nil
}

func Event(w http.ResponseWriter, r *http.Request) {
	if err := returnEvent(w, r); err != nil {
		Error(w, err)
	}
	// Just need to write a response, http.StatusOK will be
	// send automatically.
}

func main() {
	vw := &vonage.VoiceWebhook{
		Answer: http.HandlerFunc(Answer),
		Event:  http.HandlerFunc(Event),
	}
	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: vonage.NewVoiceWebhookMux(vw),
	}
	log.Printf("server listening on addr: %v", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("*** listener error: %v", err)
	}
}
