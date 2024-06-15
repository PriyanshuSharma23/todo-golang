package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/todos", app.handleGetAllTodos)
	router.Get("/todos/{id}", app.handleGetTodo)
	router.Post("/todos", app.handlePostTodos)
	router.Patch("/todos", app.handleUpdateTodos)
	router.Delete("/todos/{id}", app.handleDeleteTodo)

	return router
}
