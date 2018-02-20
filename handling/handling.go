package handling

import (
	"net/http"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/models"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// Handler encompases all request handling
type Handler struct {
	Router     http.Handler
	Config     *config.Config
	Datasource *models.Datasource
}

// NewHandler builds the handler interface and routes
func NewHandler(config *config.Config) (*Handler, error) {
	handler := new(Handler)
	handler.Config = config

	router := mux.NewRouter()

	addUserRoutes(router, handler)
	addChangeRequestRoutes(router, handler)

	// This is a "catch-all" that serves static files and logs
	// any 404s from bad requests
	router.
		PathPrefix("/").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).Then(
			http.StripPrefix("/", http.FileServer(http.Dir("static"))),
		))

	// Use gorilla's recovery handler to continue running in case of a panic
	handler.Router = handlers.RecoveryHandler()(router)

	return handler, nil
}