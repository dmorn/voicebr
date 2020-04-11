package vonage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type VoiceAnswerRequest struct {
	To               string `json:"to"`
	From             string `json:"from"`
	UUID             string `json:"uuid"`
	ConversationUUID string `json:"conversation_uuid"`
}

type VoiceEventRequest struct {
	To               string    `json:"to"`
	From             string    `json:"from"`
	UUID             string    `json:"uuid"`
	ConversationUUID string    `json:"conversation_uuid"`
	Status           string    `json:"status"`
	Direction        string    `json:"direction"`
	Timestamp        time.Time `json:"timestamp"`
}

type VoiceWebhook struct {
	// Anwer gets called when Vonage receives a call.
	Answer http.Handler
	// Event gets called each time a Vonage call changes
	// state (e.g. ringing, answered), as well as when
	// an error is occurred (e.g. we return a broken NCCO).
	Event http.Handler
}

func discardHandler(w http.ResponseWriter, r http.Request) {
	defer r.Body.Close()
	io.Copy(ioutil.Discard, r.Body)
	w.WriteHeader(http.StatusOK)
}

// NewWebhook returns a VoiceWebhook instance with default handlers,
// which do nothing but trashing the request body.
func NewWebook() *VoiceWebhook {
	return &VoiceWebhook{
		Answer: discardHandler,
		Event:  discardHandler,
	}
}

func NewVoiceWebhooksMux(hook *VoiceWebhook) *http.ServeMux {
	m := http.NewServeMux()
	m.Handle("/voice/answer", h.Answer)
	m.Handle("/voice/event", h.Event)
	return m
}
