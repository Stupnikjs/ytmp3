package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	Port int
}

func (app *Application) Routes() http.Handler {

	mux := chi.NewRouter()

	// register routes
	mux.Get("/", app.RenderAccueil)
	mux.Post("/video/id", app.PostVideoId)
	mux.Get("/download/{name}", app.DowloadSound)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}
