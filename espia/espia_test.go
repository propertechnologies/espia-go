package espia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEspiaSetup(t *testing.T) {
	espiaConfig := EspiaSetup{
		Source:      "test-product",
		AutoSession: true,
		Enabled:     true,
	}

	Espia(espiaConfig)

	if source != "test-product" {
		t.Errorf("expected source to be 'test-product', got '%s'", source)
	}

	if !enabled {
		t.Errorf("expected espia to be enabled, but it's disabled")
	}

	if session == "" {
		t.Errorf("expected auto-generated session, but got empty")
	}
}

func TestSetSession(t *testing.T) {
	testSession := "test-session-12345"
	SetSession(testSession)

	if session != testSession {
		t.Errorf("expected session to be '%s', got '%s'", testSession, session)
	}
}

func TestSetPermanentMetadata(t *testing.T) {
	SetPermanentMetadata("user_id", 12345)

	if metadata["user_id"] != 12345 {
		t.Errorf("expected user_id metadata to be 12345, got %v", metadata["user_id"])
	}
}

func TestTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	EventURL = server.URL + "/events"

	Espia(EspiaSetup{
		Source:      "test-product",
		AutoSession: true,
		Enabled:     true,
	})

	SetPermanentMetadata("user_id", 12345)

	err := Track("button_click", Metadata{
		"button_name": "submit",
	}, nil)

	if err != nil {
		t.Errorf("unexpected error while tracking event: %v", err)
	}
}

func TestTrackWithDisabled(t *testing.T) {
	Espia(EspiaSetup{
		Source:  "test-product",
		Enabled: false,
	})

	err := Track("button_click", nil, nil)

	if err != nil {
		t.Errorf("expected no error when Espia is disabled, but got %v", err)
	}
}

func TestMakeSession(t *testing.T) {
	session := makeSession()

	if len(session) == 0 {
		t.Errorf("expected non-empty session string")
	}

	if len(session) < 16 {
		t.Errorf("expected session length to be at least 16, got %d", len(session))
	}
}
