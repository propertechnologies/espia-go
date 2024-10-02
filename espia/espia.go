package espia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type (
	// Metadata is a key-value map of event data.
	Metadata = map[string]interface{}

	// EspiaSetup defines the configuration for the Espia module.
	EspiaSetup struct {
		Source      string
		Session     string
		Enabled     bool
		AutoSession bool
	}

	// Overwrites allows partial modification of Espia setup.
	Overwrites struct {
		Source  string
		Session string
	}
)

var (
	BaseURL  = "https://analytics.proper.ai"
	EventURL = BaseURL + "/events"

	source   string
	session  string
	enabled  = true
	metadata = make(Metadata)
	mu       sync.RWMutex
)

// Commit sends the event to the server.
func commit(source, category string, metadata Metadata, session string) error {
	payload := Metadata{
		"source":     source,
		"category":   category,
		"metadata":   metadata,
		"session_id": session,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", EventURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send event, status code: %d", resp.StatusCode)
	}

	return nil
}

// makeSession generates a random session ID.
func makeSession() string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	randString := randomString(16)
	return timestamp + randString
}

// randomString generates a random alphanumeric string of the given length.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Espia initializes the Espia configuration.
func Espia(config EspiaSetup) {
	mu.Lock()
	defer mu.Unlock()
	source = config.Source
	if config.AutoSession {
		session = makeSession()
	} else {
		session = config.Session
	}
	enabled = config.Enabled
	metadata = make(Metadata)
}

// SetSession sets the session ID.
func SetSession(sess string) {
	mu.Lock()
	defer mu.Unlock()
	session = sess
}

// SetPermanentMetadata sets permanent metadata that will be included in every event.
func SetPermanentMetadata(key string, value any) {
	mu.Lock()
	defer mu.Unlock()
	metadata[key] = value
}

// Track sends an event with optional metadata and overwrites.
func Track(category string, eventMetadata Metadata, overwrites *Overwrites) error {
	if !enabled {
		return nil
	}

	mu.RLock()
	defer mu.RUnlock()

	if source == "" && (overwrites == nil || overwrites.Source == "") {
		return fmt.Errorf("source is not initialized, please call Espia() before tracking events")
	}

	eventSource := source
	if overwrites != nil && overwrites.Source != "" {
		eventSource = overwrites.Source
	}

	eventSession := session
	if overwrites != nil && overwrites.Session != "" {
		eventSession = overwrites.Session
	}

	combinedMetadata := metadataWithPermanent(eventMetadata)

	return commit(eventSource, category, combinedMetadata, eventSession)
}

// metadataWithPermanent combines permanent metadata with event-specific metadata.
func metadataWithPermanent(eventMetadata Metadata) Metadata {
	mu.RLock()
	defer mu.RUnlock()

	if len(metadata) == 0 && len(eventMetadata) == 0 {
		return nil
	}

	result := make(Metadata, len(metadata)+len(eventMetadata))
	for k, v := range metadata {
		result[k] = v
	}
	for k, v := range eventMetadata {
		result[k] = v
	}

	return result
}
