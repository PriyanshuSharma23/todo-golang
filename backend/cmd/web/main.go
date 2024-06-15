package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/PriyanshuSharma23/todo-golang/internals/data"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type application struct {
	cfg struct {
		port int
		db   struct {
			dsn string
		}
	}
	logger *zerolog.Logger
	models *data.Models
}

func main() {
	var app application
	// application config
	flag.IntVar(&app.cfg.port, "port", 4000, "port for the application")

	// database config
	flag.StringVar(&app.cfg.db.dsn, "db-dsn", "", "postgres db connection url")

	flag.Parse()

	var logger = zerolog.New(os.Stdout)
	app.logger = &logger

	db, err := connectDB(&app)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to initialize connection")
	}
	defer db.Close()

	app.models = data.NewModels(db)
	logger.Info().Msg("Database connection established")

	logger.Info().Msg(fmt.Sprintf("Server starting at port: %d", app.cfg.port))
	logger.Panic().Err(app.Serve())
}

func connectDB(app *application) (*sql.DB, error) {
	db, err := sql.Open("postgres", app.cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
