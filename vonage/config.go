package vonage

type Config struct {
	// Origin Vonage should send requests to.
	Origin string `json:"origin"`
	// Application identifier.
	AppID string `json:"app_id"`
	// AppNumber is the application linked number.
	AppNumber string `json:"app_number"`
}
