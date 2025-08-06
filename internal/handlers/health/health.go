package health

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type Status struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		Health
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Status			"ok"
//	@Failure	500	{object}	ErrorResponse	"internal server error"
//	@Router		/live [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	status := Status{"ok"}
	v, err := jsoniter.Marshal(status)
	if err != nil {
		http.Error(w, errors.Wrap(err, "marshal status").Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(v)
}
