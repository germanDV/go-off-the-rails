package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"

	"github.com/germandv/go-off-the-rails/controllers"
	"github.com/germandv/go-off-the-rails/db"
	"github.com/germandv/go-off-the-rails/db/generated"
)

const (
	HttpPort              = 8787
	HttpIdleTimeout       = 1 * time.Minute
	HttpReadHeaderTimeout = 5 * time.Second
	HttpReadTimeout       = 10 * time.Second
	HttpWriteTimeout      = 15 * time.Second
	HttpTimeout           = 14 * time.Second
	DbPath                = "./db.sqlite"
)

func main() {
	dbClient, err := sql.Open("sqlite", DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: run db migrations
	tempMigrationRun(dbClient)

	moviesRepo := db.NewMoviesRepository(generated.New(dbClient))
	moviesController := controllers.NewMoviesController(moviesRepo)

	mux := &http.ServeMux{}
	moviesController.RegisterRoutes(mux)

	port := HttpPort
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		IdleTimeout:       HttpIdleTimeout,
		ReadHeaderTimeout: HttpReadHeaderTimeout,
		ReadTimeout:       HttpReadTimeout,
		WriteTimeout:      HttpWriteTimeout,
	}

	httpServer.Handler = http.TimeoutHandler(mux, HttpTimeout, "Request timed out")

	fmt.Printf("Listening on port %d\n", port)
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func tempMigrationRun(dbClient *sql.DB) {
	_, err := dbClient.Exec(`
CREATE TABLE IF NOT EXISTS orgs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  org_id TEXT NOT NULL REFERENCES orgs(id),
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);

CREATE TABLE IF NOT EXISTS movies (
  id TEXT PRIMARY KEY,
  org_id TEXT NOT NULL REFERENCES orgs(id),
  title TEXT NOT NULL,
  rating INTEGER,
  version INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migrations applied successfully")
}
