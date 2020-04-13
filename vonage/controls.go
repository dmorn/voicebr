package vonage

import (
	"encoding/json"
	"io"
)

type TalkControl struct {
	Action string `json:"action"`
	Text   string `json:"text"`
}

func NewTalkControl(text string) *TalkControl {
	return &TalkControl{
		Action: "talk",
		Text:   text,
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
