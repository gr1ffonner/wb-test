package handler

import (
	"net/http"

	_ "wb-test/api"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	// Health
	{
		router.HandleFunc("/live", h.health.Health).Methods(http.MethodGet)
	}

	// Swagger
	{
		// Redirect /swagger to /swagger/index.html
		router.Handle("/documentation", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently)).Methods(http.MethodGet)
		// Serve Swagger UI
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	}

	return router
}
