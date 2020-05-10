package vonage

type Config struct {
	// Origin Vonage should send requests to.
	Origin string `json:"origin"`
	// Application identifier.
	AppID string `json:"app_id"`
	// AppNumber is the application linked number.
	AppNumber string `json:"app_number"`
	// When doing text-to-speech, user selectable voice to use.
	// Choose one that matches the spoken language for best results.
	// https://developer.nexmo.com/voice/voice-api/guides/text-to-speech#voice-names
	VoiceName string `json:"voice_name"`
}
