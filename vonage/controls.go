package vonage

import (
	"encoding/json"
	"io"
)

type TalkControl struct {
	Action    string `json:"action"`
	Text      string `json:"text"`
	VoiceName string `json:"voice_name"`
}

func NewTalkControl(voiceName, text string) *TalkControl {
	return &TalkControl{
		Action:    "talk",
		Text:      text,
		VoiceName: voiceName,
	}
}

type RecordControl struct {
	Action      string   `json:"action"`
	Format      string   `json:"format"`
	Timeout     int      `json:"timeout"`
	BeepStart   bool     `json:"beepStart"`
	EventURL    []string `json:"eventUrl"`
	EventMethod string   `json:"eventMethod"`
}

func NewRecordControl(url string) *RecordControl {
	return &RecordControl{
		Action:      "record",
		Format:      "mp3",
		Timeout:     30,
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
