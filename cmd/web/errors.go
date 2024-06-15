package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(_ *http.Request, err error) {
	app.logger.Error().Err(err).Msg("")
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	msg := envelope{
		"error": data,
	}

	err := app.writeJSON(w, r, status, msg)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	msg := `something went wrong while processing your request`
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, fmt.Sprint(err))
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	msg := `the requested resource cannot be found`
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) unproccessableEntityError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}
