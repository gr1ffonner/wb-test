package handler

import "wb-test/internal/handlers/health"

type Handler struct {
	health *health.Handler
}

func NewHandler() *Handler {
	return &Handler{
		health: health.NewHandler(),
	}
}
