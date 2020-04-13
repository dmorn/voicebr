package enginevonage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jecoz/voicebr/vonage"
)

func AllowBroadcastFrom(caller string) bool {
	// TODO(jecoz): load preferences and caller whitelist.
	return true
}

func Answer(w http.ResponseWriter, r *http.Request) error {
	p, err := vonage.ParseVoiceAnswerRequest(r)
	if err != nil {
		return &httpError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("answer request: %w", err),
		}
	}
	if !AllowBroadcastFrom(p.From) {
		return &httpError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("%v is not allowed to broadcast", p.From),
		}
	}

	log.Printf("*** reveived call (%v) from: %v", p.UUID, p.From)

	c := vonage.NewTalkControl("Hi from Jecoz")
	if err = vonage.NewEncoder(w).EncodeControls(c); err != nil {
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
