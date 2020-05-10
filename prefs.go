package voley

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jecoz/voley/vonage"
	"github.com/tailscale/hujson"
)

// Prefs are the user modifiable preferences
// of the voley server. It embeds also vonage's
// configuration, possibly in the future its alternative.
type Prefs struct {
	// List of callers that are allowed to initiate
	// broadcast requests.
	Broadcasters []string `json:"broadcasters"`
	// Message told to the caller before the recording starts.
	BroadcastGreetMsg string         `json:"broadcast_greet_msg"`
	ExternalOrigin    string         `json:"external_origin"`
	Port              int            `json:"port"`
	Vonage            *vonage.Config `json:"vonage"`
}

func WritePrefs(w io.Writer, p *Prefs) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(&p); err != nil {
		return fmt.Errorf("write preferences: %w", err)
	}
	return nil
}

func SavePrefs(filename string, p *Prefs) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("save preferences: %w", err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save preferences: %w", err)
	}
	defer f.Close()
	return WritePrefs(f, p)
}

func ReadPrefs(r io.Reader, p *Prefs) error {
	if err := hujson.NewDecoder(r).Decode(&p); err != nil {
		return fmt.Errorf("read preferences: %w", err)
	}
	return nil
}

func LoadPrefs(filename string, p *Prefs) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}
	defer f.Close()
	return ReadPrefs(f, p)
}
