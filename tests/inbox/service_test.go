package inbox_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"server/internal/inbox"
	"server/internal/storage"
)

func TestServiceCreatesAndListsInboxes(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	service := inbox.NewService(inbox.NewRepository(database))
	createdInbox, err := service.Create(context.Background(), "  Support  ")
	if err != nil {
		t.Fatalf("Create() returned an error: %v", err)
	}

	if createdInbox.ID != 1 {
		t.Fatalf("created inbox ID = %d, want 1", createdInbox.ID)
	}

	if createdInbox.Name != "Support" {
		t.Fatalf("created inbox name = %q, want %q", createdInbox.Name, "Support")
	}

	inboxes, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() returned an error: %v", err)
	}

	if len(inboxes) != 1 || inboxes[0].Name != "Support" {
		t.Fatalf("listed inboxes = %#v, want one Support inbox", inboxes)
	}
}

func TestServiceRejectsBlankInboxName(t *testing.T) {
	database, err := storage.OpenAt(filepath.Join(t.TempDir(), "hooklens.db"))
	if err != nil {
		t.Fatalf("OpenAt() returned an error: %v", err)
	}
	defer database.Close()

	service := inbox.NewService(inbox.NewRepository(database))
	if _, err := service.Create(context.Background(), "   "); !errors.Is(err, inbox.ErrNameRequired) {
		t.Fatalf("Create() error = %v, want %v", err, inbox.ErrNameRequired)
	}
}
