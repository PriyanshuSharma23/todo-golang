package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/PriyanshuSharma23/todo-golang/internals/data"
	"github.com/PriyanshuSharma23/todo-golang/internals/validator"
	"github.com/go-chi/chi/v5"
)

func (app *application) handleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.models.TodosModel.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, envelope{
		"data": todos,
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	todo, err := app.models.TodosModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundError(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, envelope{
		"data": todo,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}

type postTodoInput struct {
	Title string    `json:"title"`
	DueOn time.Time `json:"due_on"`
	Body  string    `json:"body"`
}

func (app *application) handlePostTodos(w http.ResponseWriter, r *http.Request) {
	var input postTodoInput
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	todo := data.Todo{
		Title:      input.Title,
		IsComplete: false,
		DueOn:      input.DueOn,
		Body:       input.Body,
	}

	var v = validator.New()
	data.VerifyTodo(v, &todo)
	if !v.Valid() {
		app.unproccessableEntityError(w, r, v.Errors)
		return
	}

	err = app.models.TodosModel.Insert(&todo)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusCreated, envelope{"data": todo})
	if err != nil {
		app.serverError(w, r, err)
	}
}

type patchTodoInput struct {
	ID         *int       `json:"id"`
	Title      *string    `json:"title"`
	IsComplete *bool      `json:"is_complete"`
	DueOn      *time.Time `json:"due_on"`
	Body       *string    `json:"body"`
}

func (app *application) handleUpdateTodos(w http.ResponseWriter, r *http.Request) {
	var input patchTodoInput
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if input.ID == nil {
		var v = validator.New()
		v.AddError("id", "id required")
		app.unproccessableEntityError(w, r, v.Errors)
		return
	}

	todo, err := app.models.TodosModel.Get(*input.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundError(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	if input.Title != nil {
		todo.Title = *input.Title
	}

	if input.Body != nil {
		todo.Body = *input.Body
	}

	if input.IsComplete != nil {
		todo.IsComplete = *input.IsComplete
	}

	if input.DueOn != nil {
		todo.DueOn = *input.DueOn
	}

	var v = validator.New()
	data.VerifyTodo(v, todo)
	if !v.Valid() {
		app.unproccessableEntityError(w, r, v.Errors)
		return
	}

	err = app.models.TodosModel.Update(todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.editConflictResponse(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	if err = app.writeJSON(w, r, http.StatusOK, envelope{"data": todo}); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	todo, err := app.models.TodosModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundError(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, envelope{
		"data": todo,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
