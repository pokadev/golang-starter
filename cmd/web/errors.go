package main

import (
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) reportError(err error) {
	trace := debug.Stack()

	app.logger.Error(err, trace)

	if app.config.notifications.email != "" {
		data := app.newEmailData()
		data["Timestamp"] = time.Now().UTC().Format(time.RFC3339)
		data["Error"] = err.Error()
		data["Trace"] = string(trace)

		err := app.mailer.Send(app.config.notifications.email, data, "error-notification.tmpl")
		if err != nil {
			trace := debug.Stack()
			app.logger.Error(err, trace)
		}
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportError(err)

	message := "The server encountered a problem and could not process your request"
	http.Error(w, message, http.StatusInternalServerError)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	http.Error(w, message, http.StatusNotFound)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	message := "You must be authenticated to access this resource"
	http.Error(w, message, http.StatusUnauthorized)
}
