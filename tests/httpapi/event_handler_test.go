package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"server/internal/health"
	"server/internal/httpapi"
	"server/internal/inbox"
	"server/internal/ping"
	"server/internal/storage"
)

func TestEventRoutesCreateAndList(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	inboxService := inbox.NewService(inbox.NewRepository(database))
	createdInbox, err := inboxService.Create(t.Context(), "Support")
	if err != nil {
		t.Fatalf("Create() returned an error: %v", err)
	}

	router := httpapi.NewRouter(
		httpapi.NewHealthHandler(health.NewService("hooklens", "test")),
		httpapi.NewPingHandler(ping.NewService()),
		httpapi.NewInboxHandler(inboxService),
		httpapi.NewEventHandler(inboxService),
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/inboxes/1/events",
		strings.NewReader(`{"type":"order.created"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-ID", "request-123")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("POST event status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}

	var createdResponse struct {
		Data inbox.InboxEvent `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &createdResponse); err != nil {
		t.Fatalf("decode POST response: %v", err)
	}
	if createdResponse.Data.InboxID != createdInbox.ID {
		t.Fatalf("created event inbox ID = %d, want %d", createdResponse.Data.InboxID, createdInbox.ID)
	}

	request = httptest.NewRequest(http.MethodGet, "/api/inboxes/1/events", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("GET events status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}

	var listedResponse struct {
		Data []inbox.InboxEvent `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &listedResponse); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if len(listedResponse.Data) != 1 || listedResponse.Data[0].Body != `{"type":"order.created"}` {
		t.Fatalf("listed events = %#v, want one event with original body", listedResponse.Data)
	}
}

func TestEventRoutesReturnNotFoundForMissingInbox(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	inboxService := inbox.NewService(inbox.NewRepository(database))
	router := httpapi.NewRouter(
		httpapi.NewHealthHandler(health.NewService("hooklens", "test")),
		httpapi.NewPingHandler(ping.NewService()),
		httpapi.NewInboxHandler(inboxService),
		httpapi.NewEventHandler(inboxService),
	)

	request := httptest.NewRequest(http.MethodGet, "/api/inboxes/999/events", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("missing inbox status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
