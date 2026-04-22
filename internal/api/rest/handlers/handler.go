package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type (
	baseHandler struct {
		log *slog.Logger
	}
)

func (h *baseHandler) respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.log.Error("failed to encode response", "error", err)
		}
	}
}

func (h *baseHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respond(w, status, map[string]string{"error": message})
}
