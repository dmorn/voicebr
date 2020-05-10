package enginevonage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jecoz/voiley"
	"github.com/jecoz/voiley/vonage"
)

// RecordHandler is an "answer" webhook which tells Vonage to record the call.
type RecordHandler struct {
	// When the record is completed, RecordCallbackURL is contacted with
	// record data.
	RecordCallbackURL string

	*vonage.Config
	*voiley.Prefs
}

func AllowBroadcastFrom(from string, broadcasters []string) bool {
	for _, v := range broadcasters {
		if v == from {
			return true
		}
	}
	return false
}

func (a *RecordHandler) ServeHTTPReturn(w http.ResponseWriter, r *http.Request) error {
	p, err := vonage.ParseVoiceAnswerRequest(r)
	if err != nil {
		return &httpError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("answer request: %w", err),
		}
	}
	if !AllowBroadcastFrom(p.From, a.Broadcasters) {
		return &httpError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("%v is not allowed to broadcast", p.From),
		}
	}

	log.Printf("*** reveived call (%v) from: %v", p.UUID, p.From)

	nccos := []interface{}{}
	if msg := a.BroadcastGreetMsg; msg != "" {
		nccos = append(nccos, vonage.NewTalkControl(a.VoiceName, msg))
	}
	nccos = append(nccos, vonage.NewRecordControl(a.RecordCallbackURL))

	mw := io.MultiWriter(w, os.Stderr)
	if err = vonage.NewEncoder(mw).EncodeControls(nccos...); err != nil {
		return fmt.Errorf("answer response: %w", err)
	}
	return nil
}

func NewRecordHandlerFunc(callback string, c *vonage.Config, p *voiley.Prefs) ReturnHandlerFunc {
	a := &RecordHandler{callback, c, p}
	return a.ServeHTTPReturn
}

func Event(w http.ResponseWriter, r *http.Request) error {
	_, err := vonage.ParseVoiceEventRequest(r)
	if err != nil {
		return &httpError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("answer request: %w", err),
		}
	}

	// log.Printf("event: %+v", p)
	return nil
}
