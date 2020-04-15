package vonage

import (
	"os"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/tailscale/hujson"
)

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

func WriteConfig(w io.Writer, c *Config) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(&c); err != nil {
		return fmt.Errorf("write preferences: %w", err)
	}
	return nil
}

func SaveConfig(filename string, c *Config) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("save preferences: %w", err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save preferences: %w", err)
	}
	defer f.Close()
	return WriteConfig(f, c)
}

func ReadConfig(r io.Reader, c *Config) error {
	if err := hujson.NewDecoder(r).Decode(&c); err != nil {
		return fmt.Errorf("read preferences: %w", err)
	}
	return nil
}

func LoadConfig(filename string, c *Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("load preferences: %w", err)
	}
	defer f.Close()
	return ReadConfig(f, c)
}
