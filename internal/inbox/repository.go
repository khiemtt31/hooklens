package inbox

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (r *Repository) Get(ctx context.Context, id int64) (Inbox, error) {
	var (
		inbox     Inbox
		createdAt string
	)

	err := r.database.QueryRowContext(
		ctx,
		"SELECT id, name, created_at FROM inboxes WHERE id = ?",
		id,
	).Scan(&inbox.ID, &inbox.Name, &createdAt)
	if err != nil {
		return Inbox{}, fmt.Errorf("get inbox: %w", err)
	}

	inbox.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return Inbox{}, fmt.Errorf("parse inbox created time: %w", err)
	}

	return inbox, nil
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

func (r *Repository) CreateEvent(
	ctx context.Context,
	inboxID int64,
	method string,
	headers map[string][]string,
	body string,
) (InboxEvent, error) {
	receivedAt := time.Now().UTC()
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return InboxEvent{}, fmt.Errorf("encode event headers: %w", err)
	}

	result, err := r.database.ExecContext(
		ctx,
		"INSERT INTO inbox_events (inbox_id, method, headers, body, received_at) VALUES (?, ?, ?, ?, ?)",
		inboxID,
		method,
		headersJSON,
		body,
		receivedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return InboxEvent{}, fmt.Errorf("insert inbox event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return InboxEvent{}, fmt.Errorf("get inbox event id: %w", err)
	}

	return InboxEvent{
		ID:         id,
		InboxID:    inboxID,
		Method:     method,
		Headers:    headers,
		Body:       body,
		ReceivedAt: receivedAt,
	}, nil
}

func (r *Repository) ListEvents(ctx context.Context, inboxID int64) ([]InboxEvent, error) {
	rows, err := r.database.QueryContext(
		ctx,
		"SELECT id, inbox_id, method, headers, body, received_at FROM inbox_events WHERE inbox_id = ? ORDER BY id ASC",
		inboxID,
	)
	if err != nil {
		return nil, fmt.Errorf("select inbox events: %w", err)
	}
	defer rows.Close()

	events := make([]InboxEvent, 0)
	for rows.Next() {
		var (
			event       InboxEvent
			headersJSON string
			receivedAt  string
		)

		if err := rows.Scan(
			&event.ID,
			&event.InboxID,
			&event.Method,
			&headersJSON,
			&event.Body,
			&receivedAt,
		); err != nil {
			return nil, fmt.Errorf("scan inbox event: %w", err)
		}

		event.Headers = make(map[string][]string)
		if err := json.Unmarshal([]byte(headersJSON), &event.Headers); err != nil {
			return nil, fmt.Errorf("decode inbox event headers: %w", err)
		}

		event.ReceivedAt, err = time.Parse(time.RFC3339Nano, receivedAt)
		if err != nil {
			return nil, fmt.Errorf("parse inbox event received time: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inbox events: %w", err)
	}

	return events, nil
}
