package vonage

import (
	"encoding/json"
	"io"
)

type TalkControl struct {
	Action    string `json:"action,omitempty"`
	Text      string `json:"text,omitempty"`
	VoiceName string `json:"voice_name,omitempty"`
}

func NewTalkControl(voiceName, text string) *TalkControl {
	return &TalkControl{
		Action:    "talk",
		Text:      text,
		VoiceName: voiceName,
	}
}

type RecordControl struct {
	Action      string   `json:"action,omitempty"`
	Format      string   `json:"format,omitempty"`
	Timeout     int      `json:"timeout,omitempty"`
	BeepStart   bool     `json:"beepStart,omitempty"`
	EventURL    []string `json:"eventUrl,omitempty"`
	EventMethod string   `json:"eventMethod,omitempty"`
}

func NewRecordControl(url string) *RecordControl {
	return &RecordControl{
		Action: "record",
		Format: "mp3",
		//Timeout:     30,
		BeepStart:   true,
		EventURL:    []string{url},
		EventMethod: "POST",
	}
}

type Encoder struct {
	e *json.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{json.NewEncoder(w)}
}

func (e *Encoder) EncodeControls(items ...interface{}) error {
	return e.e.Encode(items)
}
