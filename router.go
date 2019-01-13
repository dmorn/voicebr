package voicebr

import (
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	"io"

	"github.com/gorilla/mux"
)

var (
	HostAddr = "http://localhost:4001"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/record/voice/answer", RecordAnswerHandler)
	r.HandleFunc("/record/voice/event", RecordEventHandler)
	r.HandleFunc("/store/recording/event", StoreRecordingEventHandler)
	r.Use(loggingMiddleware)

	return r
}

func RecordAnswerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		map[string]interface{}{
			"action": "talk",
			"text": "Hello",
		},
		map[string]interface{}{
			"action": "record",
			"beepStart": true,
			"eventUrl": HostAddr+"/store/recording/event",
			"endOnKey": 1,
		},
	})
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

func StoreRecordingEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	defer r.Body.Close()

	var content struct{
		ConversationUUID string `json:"conversation_uuid"`
		RecordingURL string `json:"recording_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Download mp3 file with the recording. It will
	// later be used into the outbound calls.
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
