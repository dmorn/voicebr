package vonage

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
