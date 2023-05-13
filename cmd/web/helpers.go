package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Logger) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError),
	http.StatusInternalServerError)
	}

func (app *Logger) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Logger) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
