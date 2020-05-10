package vonage

import (
	"encoding/json"
	"fmt"
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

// Depending on the request Method, either decodes and closes the request Body
// or parses the request Form.
func ParseVoiceAnswerRequest(r *http.Request) (*VoiceAnswerRequest, error) {
	switch r.Method {
	case "POST":
		var p *VoiceAnswerRequest
		if err := unmarshalRequest(&p, r); err != nil {
			return nil, fmt.Errorf("parse voice answer: %w", err)
		}
		return p, nil
	case "GET":
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("parse voice answer: %w", err)
		}
		return &VoiceAnswerRequest{
			To:               r.Form.Get("to"),
			From:             r.Form.Get("from"),
			UUID:             r.Form.Get("uuid"),
			ConversationUUID: r.Form.Get("conversation_uuid"),
		}, nil
	default:
		return nil, fmt.Errorf("parse voice answer: unsupported request method")
	}
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

// Depending on the request Method, either decodes and closes the request Body
// or parses the request Form.
func ParseVoiceEventRequest(r *http.Request) (*VoiceEventRequest, error) {
	switch r.Method {
	case "POST":
		var p *VoiceEventRequest
		if err := unmarshalRequest(&p, r); err != nil {
			return nil, fmt.Errorf("parse voice event: %w", err)
		}
		return p, nil
	case "GET":
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("parse voice event: %w", err)
		}
		t, err := time.Parse(time.RFC3339, r.Form.Get("timestamp"))
		if err != nil {
			return nil, fmt.Errorf("parse voice event: %w", err)
		}
		return &VoiceEventRequest{
			To:               r.Form.Get("to"),
			From:             r.Form.Get("from"),
			UUID:             r.Form.Get("uuid"),
			ConversationUUID: r.Form.Get("conversation_uuid"),
			Status:           r.Form.Get("status"),
			Direction:        r.Form.Get("direction"),
			Timestamp:        t,
		}, nil
	default:
		return nil, fmt.Errorf("parse voice answer: unsupported request method")
	}
}

type VoiceWebhook struct {
	// Anwer gets called when Vonage receives a call.
	Answer http.Handler
	// Event gets called each time a Vonage call changes
	// state (e.g. ringing, answered), as well as when
	// an error is occurred (e.g. we return a broken NCCO).
	Event http.Handler
}

// NewWebhook returns a VoiceWebhook instance with default handlers,
// which do nothing but trashing the request body.
func NewWebook() *VoiceWebhook {
	return &VoiceWebhook{
		Answer: http.HandlerFunc(discardHandler),
		Event:  http.HandlerFunc(discardHandler),
	}
}

func NewVoiceWebhookMux(hook *VoiceWebhook) *http.ServeMux {
	m := http.NewServeMux()
	m.Handle("/answers", hook.Answer)
	m.Handle("/events", hook.Event)
	return m
}

func discardHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	io.Copy(ioutil.Discard, r.Body)
	w.WriteHeader(http.StatusOK)
}

func unmarshalRequest(i interface{}, r *http.Request) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		// Discard the remaining bytes.
		io.Copy(ioutil.Discard, r.Body)
		return fmt.Errorf("unmarshal request: %w", err)
	}
	return nil
}
