package health

import (
	"net/http"

	httputils "wb-test/pkg/utils/http-utils"
)

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
//	@Success	200	{object}	httputils.Status		"ok"
//	@Failure	500	{object}	httputils.ErrorResponse	"internal server error"
//	@Router		/live [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httputils.WriteResponse(w, http.StatusOK, "ok", nil, nil)
}
