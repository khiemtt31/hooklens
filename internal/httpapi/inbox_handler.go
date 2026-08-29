package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"server/internal/inbox"
)

type InboxHandler struct {
	service *inbox.Service
}

func NewInboxHandler(service *inbox.Service) *InboxHandler {
	return &InboxHandler{service: service}
}

type createInboxRequest struct {
	Name string `json:"name"`
}

func (h *InboxHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request createInboxRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		GlobalErrorHandler(w, http.StatusBadRequest, "invalid request body")

		return
	}

	createdInbox, err := h.service.Create(r.Context(), request.Name)
	if err != nil {
		if errors.Is(err, inbox.ErrNameRequired) {
			GlobalErrorHandler(w, http.StatusBadRequest, err.Error())

			return
		}

		GlobalErrorHandler(w, http.StatusInternalServerError, "failed to create inbox")

		return
	}

	writeResponse(w, http.StatusCreated, "inbox created", createdInbox)
}

func (h *InboxHandler) Get(w http.ResponseWriter, r *http.Request) {
	inboxes, err := h.service.List(r.Context())
	if err != nil {
		GlobalErrorHandler(w, http.StatusInternalServerError, "failed to list inboxes")

		return
	}

	writeResponse(w, http.StatusOK, "inboxes listed", inboxes)
}
