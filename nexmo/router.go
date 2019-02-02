/// Broadcast voice messages to a set of recipients.
/// Copyright (C) 2019 Daniel Morandini (jecoz)
///
/// This program is free software: you can redistribute it and/or modify
/// it under the terms of the GNU General Public License as published by
/// the Free Software Foundation, either version 3 of the License, or
/// (at your option) any later version.
///
/// This program is distributed in the hope that it will be useful,
/// but WITHOUT ANY WARRANTY; without even the implied warranty of
/// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
/// GNU General Public License for more details.
///
/// You should have received a copy of the GNU General Public License
/// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package nexmo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const recFormat = "mp3"

type Storage interface {
	ContactsProvider
	RecFileHandler() http.Handler
	WriteRec(src io.Reader, fileName string) (string, error)
}

func NewRouter(c *Client, s Storage, hostAddr string) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/record/voice/answer", makeRecordAnswerHandler(hostAddr))
	r.HandleFunc("/record/voice/event", LogEventHandler)
	r.HandleFunc("/store/recording/event", makeStoreRecordingEventHandler(s, c))
	r.HandleFunc("/play/recording/event", LogEventHandler)
	r.HandleFunc("/play/recording/{name}", makePlayRecordingHandler(hostAddr))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", s.RecFileHandler()))
	r.Use(loggingMiddleware)

	return r
}

func makeRecordAnswerHandler(hostAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"action": "talk",
				"text":   "Talk now",
			},
			{
				"action":    "record",
				"beepStart": true,
				"format":    recFormat,
				"eventUrl":  []string{hostAddr + "/store/recording/event"},
				"endOnKey":  1,
			},
		})
	}
}

func LogEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	defer func() {
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
	}()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		log.Printf("Error: unable to read body: %v", err)
	}

	log.Printf("Event received: %v", buf.String())
}

func makeStoreRecordingEventHandler(s Storage, c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		defer r.Body.Close()

		var content struct {
			ConversationUUID string `json:"conversation_uuid"`
			RecordingURL     string `json:"recording_url"`
			RecordingUUID    string `json:"recording_uuid"`
		}
		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			log.Printf("Error: unable to decode recorinding event: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Download mp3 file with the recording. It will
		// later be used into the outbound calls.
		resp, err := c.Get(content.RecordingURL)
		if err != nil {
			log.Printf("RecordingEventHandler error: unable to download file: %v", err)
			return
		}
		defer resp.Body.Close()

		recName := content.RecordingUUID + "." + recFormat
		if _, err = s.WriteRec(resp.Body, recName); err != nil {
			log.Println(err)
			return
		}

		// Make outbound phone call that will play the saved
		// recording.
		c.Call(s, recName)
	}
}

func makePlayRecordingHandler(hostAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"action": "talk",
				"text":   "Recorded message",
			},
			{
				"action":    "stream",
				"streamUrl": []string{hostAddr + "/static/" + name},
			},
		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[%s] %s", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
