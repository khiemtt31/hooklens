package inbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{database: database}
}

func (r *Repository) Create(ctx context.Context, name string) (Inbox, error) {
	createdAt := time.Now().UTC()

	result, err := r.database.ExecContext(
		ctx,
		"INSERT INTO inboxes (name, created_at) VALUES (?, ?)",
		name,
		createdAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return Inbox{}, fmt.Errorf("insert inbox: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Inbox{}, fmt.Errorf("get inbox id: %w", err)
	}

	return Inbox{ID: id, Name: name, CreatedAt: createdAt}, nil
}

func (r *Repository) List(ctx context.Context) ([]Inbox, error) {
	rows, err := r.database.QueryContext(
		ctx,
		"SELECT id, name, created_at FROM inboxes ORDER BY id ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("select inboxes: %w", err)
	}
	defer rows.Close()

	inboxes := make([]Inbox, 0)
	for rows.Next() {
		var (
			inbox     Inbox
			createdAt string
		)

		if err := rows.Scan(&inbox.ID, &inbox.Name, &createdAt); err != nil {
			return nil, fmt.Errorf("scan inbox: %w", err)
		}

		inbox.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse inbox created time: %w", err)
		}

		inboxes = append(inboxes, inbox)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inboxes: %w", err)
	}

	return inboxes, nil
}
