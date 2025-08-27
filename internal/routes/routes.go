package routes

import (
	"github.com/augustoapg/censysKvStore/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/kv", func(r chi.Router) {
		r.Get("/{key}", app.KVHandler.HandleGetKvByKey)
		r.Post("/", app.KVHandler.HandleSetKv)
	})

	return r
}
