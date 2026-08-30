package httpapi

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"server/internal/inbox"
)

const maxInboxEventBodyBytes = 1 << 20

var errInvalidInboxID = errors.New("invalid inbox id")

// EventHandler captures and lists HTTP events for an inbox.
type EventHandler struct {
	service *inbox.Service
}

func NewEventHandler(service *inbox.Service) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	inboxID, err := inboxIDFromRequest(r)
	if err != nil {
		GlobalErrorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxInboxEventBodyBytes))
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			GlobalErrorHandler(w, http.StatusRequestEntityTooLarge, "event body is too large")
			return
		}

		GlobalErrorHandler(w, http.StatusBadRequest, "failed to read event body")
		return
	}

	createdEvent, err := h.service.CreateEvent(
		r.Context(),
		inboxID,
		r.Method,
		map[string][]string(r.Header),
		string(body),
	)
	if err != nil {
		if errors.Is(err, inbox.ErrInboxNotFound) {
			GlobalErrorHandler(w, http.StatusNotFound, err.Error())
			return
		}

		GlobalErrorHandler(w, http.StatusInternalServerError, "failed to create inbox event")
		return
	}

	writeResponse(w, http.StatusCreated, "inbox event created", createdEvent)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	inboxID, err := inboxIDFromRequest(r)
	if err != nil {
		GlobalErrorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.service.ListEvents(r.Context(), inboxID)
	if err != nil {
		if errors.Is(err, inbox.ErrInboxNotFound) {
			GlobalErrorHandler(w, http.StatusNotFound, err.Error())
			return
		}

		GlobalErrorHandler(w, http.StatusInternalServerError, "failed to list inbox events")
		return
	}

	writeResponse(w, http.StatusOK, "inbox events listed", events)
}

func inboxIDFromRequest(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errInvalidInboxID
	}

	return id, nil
}
