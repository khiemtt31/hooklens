package inbox_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"server/internal/inbox"
	"server/internal/storage"
)

func TestServiceCreatesAndListsEvents(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	service := inbox.NewService(inbox.NewRepository(database))
	createdInbox, err := service.Create(context.Background(), "Support")
	if err != nil {
		t.Fatalf("Create() returned an error: %v", err)
	}

	wantHeaders := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Request-ID": {"request-123"},
	}
	wantBody := `{"type":"order.created","orderId":"123"}`

	createdEvent, err := service.CreateEvent(
		context.Background(),
		createdInbox.ID,
		"POST",
		wantHeaders,
		wantBody,
	)
	if err != nil {
		t.Fatalf("CreateEvent() returned an error: %v", err)
	}

	if createdEvent.ID != 1 {
		t.Fatalf("created event ID = %d, want 1", createdEvent.ID)
	}
	if createdEvent.InboxID != createdInbox.ID {
		t.Fatalf("created event inbox ID = %d, want %d", createdEvent.InboxID, createdInbox.ID)
	}
	if createdEvent.Method != "POST" || createdEvent.Body != wantBody {
		t.Fatalf("created event = %#v, want POST and original body", createdEvent)
	}
	if createdEvent.ReceivedAt.IsZero() {
		t.Fatal("created event timestamp is zero")
	}

	events, err := service.ListEvents(context.Background(), createdInbox.ID)
	if err != nil {
		t.Fatalf("ListEvents() returned an error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("listed events length = %d, want 1", len(events))
	}
	if !reflect.DeepEqual(events[0].Headers, wantHeaders) {
		t.Fatalf("listed event headers = %#v, want %#v", events[0].Headers, wantHeaders)
	}
	if events[0].Body != wantBody {
		t.Fatalf("listed event body = %q, want %q", events[0].Body, wantBody)
	}
}

func TestServiceRejectsEventsForMissingInbox(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	service := inbox.NewService(inbox.NewRepository(database))
	_, err = service.CreateEvent(context.Background(), 999, "POST", nil, "body")
	if !errors.Is(err, inbox.ErrInboxNotFound) {
		t.Fatalf("CreateEvent() error = %v, want %v", err, inbox.ErrInboxNotFound)
	}
}
