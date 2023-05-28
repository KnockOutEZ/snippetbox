package main

import (
	"net/http"
)

func (app *Application) routes(mux *http.ServeMux) http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	return app.logRequest(secureHeaders(mux))
}
