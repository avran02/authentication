package controller

import (
	"log/slog"
	"net/http"
)

func apiError(w http.ResponseWriter, status int, err error) {
	slog.Error("failed to unmarshal JSON", "error", err.Error())
	w.WriteHeader(status)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		slog.Error("failed to write response", "error", err.Error)
	}
}
