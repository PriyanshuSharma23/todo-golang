package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PriyanshuSharma23/todo-golang/internals/validator"
)

type TodosModel struct {
	db *sql.DB
}

type Todo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	IsComplete bool      `json:"is_complete"`
	DueOn      time.Time `json:"due_on"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	Version    int       `json:"-"`
}

func VerifyTodo(v *validator.Validator, todo *Todo) {
	v.Check(validator.MinChars(todo.Title, 3), "title", "must have ateast 3 characters")
	v.Check(validator.MaxChars(todo.Title, 50), "title", "must have atmost 50 characters")
	v.Check(validator.MaxChars(todo.Body, 200), "body", "must have atmost 200 characters")
}

func (m *TodosModel) GetAll() ([]Todo, error) {
	stmt := `SELECT id, title, is_complete, due_on, body, created_at, version FROM todos`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var todos []Todo = make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.IsComplete,
			&todo.DueOn,
			&todo.Body,
			&todo.CreatedAt,
			&todo.Version,
		)

		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return todos, nil
}

func (m *TodosModel) Get(id int) (*Todo, error) {
	var todo Todo

	stmt := `SELECT id, title, is_complete, due_on, body, created_at, version FROM todos WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, stmt, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.IsComplete,
		&todo.DueOn,
		&todo.Body,
		&todo.CreatedAt,
		&todo.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func (m *TodosModel) Insert(todo *Todo) error {
	stmt := `INSERT INTO todos (title, is_complete, due_on, body) VALUES ($1, $2, $3, $4) RETURNING id, version, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{todo.Title, todo.IsComplete, todo.DueOn, todo.Body}

	err := m.db.QueryRowContext(ctx, stmt, args...).Scan(&todo.ID, &todo.Version, &todo.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *TodosModel) Update(todo *Todo) error {
	stmt := `
			UPDATE todos SET title=$1, is_complete=$2, due_on=$3, body=$4, version=version+1
			WHERE id=$5 AND version=$6
			returning version
			`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{todo.Title, todo.IsComplete, todo.DueOn, todo.Body, todo.ID, todo.Version}
	err := m.db.QueryRowContext(ctx, stmt, args...).Scan(&todo.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m *TodosModel) Delete(id int) (*Todo, error) {
	var todo Todo
	stmt := `
			DELETE FROM todos WHERE id = $1 
			RETURNING id, title, is_complete, due_on, body, created_at, version
			`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, stmt, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.IsComplete,
		&todo.DueOn,
		&todo.Body,
		&todo.CreatedAt,
		&todo.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &todo, nil

}
