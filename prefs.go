package callrelay

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tailscale/hujson"
)

// Prefs are the user modifiable preferences
// of the callrelay server.
type Prefs struct {
	// List of callers that are allowed to initiate
	// broadcast requests.
	Broadcasters []string
	// Message told to the caller before the recording starts.
	BroadcastGreetMsg string
	ExternalOrigin    string
	Port              int
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
