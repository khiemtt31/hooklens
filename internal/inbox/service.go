package inbox

import (
	"context"
	"errors"
	"strings"
)

var ErrNameRequired = errors.New("inbox name is required")

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

func (s *Service) List(ctx context.Context) ([]Inbox, error) {
	return s.repository.List(ctx)
}
