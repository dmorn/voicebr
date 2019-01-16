package voicebr

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(rootDir, hostAddr string) *mux.Router {
	s := &Store{
		RootDir: rootDir,
	}

	r := mux.NewRouter()
	r.HandleFunc("/record/voice/answer", makeRecordAnswerHandler(hostAddr))
	r.HandleFunc("/record/voice/event", RecordEventHandler)
	r.HandleFunc("/store/recording/event", makeStoreRecordingEventHandler(s))
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
				"eventUrl":  []string{hostAddr + "/store/recording/event"},
				"endOnKey":  1,
			},
		})
	}
}

func RecordEventHandler(w http.ResponseWriter, r *http.Request) {
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

func makeStoreRecordingEventHandler(s *Store) http.HandlerFunc {
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

		resp, err := http.DefaultClient.Get(content.RecordingURL)
		if err != nil {
			log.Printf("RecordingEventHandler error: unable to download file: %v", err)
			return
		}
		defer resp.Body.Close()

		if err = s.PutRec(resp.Body, content.RecordingUUID); err != nil {
			log.Println(err)
			return
		}

		// TODO: Make outbound phone call that will play the saved
		// recording.
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
