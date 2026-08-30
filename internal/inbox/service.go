package inbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNameRequired  = errors.New("inbox name is required")
	ErrInboxNotFound = errors.New("inbox not found")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, name string) (Inbox, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Inbox{}, ErrNameRequired
	}

	return s.repository.Create(ctx, name)
}

func (s *Service) Get(ctx context.Context, id int64) (Inbox, error) {
	inbox, err := s.repository.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Inbox{}, ErrInboxNotFound
	}
	if err != nil {
		return Inbox{}, fmt.Errorf("get inbox: %w", err)
	}

	return inbox, nil
}

func (s *Service) List(ctx context.Context) ([]Inbox, error) {
	return s.repository.List(ctx)
}

func (s *Service) CreateEvent(
	ctx context.Context,
	inboxID int64,
	method string,
	headers map[string][]string,
	body string,
) (InboxEvent, error) {
	if _, err := s.Get(ctx, inboxID); err != nil {
		return InboxEvent{}, err
	}

	return s.repository.CreateEvent(ctx, inboxID, method, headers, body)
}

func (s *Service) ListEvents(ctx context.Context, inboxID int64) ([]InboxEvent, error) {
	if _, err := s.Get(ctx, inboxID); err != nil {
		return nil, err
	}

	return s.repository.ListEvents(ctx, inboxID)
}
