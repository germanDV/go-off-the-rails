package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"

	"github.com/germandv/go-off-the-rails/controllers"
	"github.com/germandv/go-off-the-rails/domain"
)

const (
	HttpPort = 8787
	DbPath   = "./db.sqlite"
)

func main() {
	dbClient, err := sql.Open("sqlite", DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: run db migrations
	tempMigrationRun(dbClient)

	// TODO: read keys from env
	tokenizer, err := domain.NewTokenizer(tempGetPrivateKey(), tempGetPublicKey())
	if err != nil {
		log.Fatal(err)
	}

	server := controllers.NewServer(controllers.ServerConfig{
		Port:      HttpPort,
		DBClient:  dbClient,
		Tokenizer: tokenizer,
	})

	fmt.Printf("Listening on port %d\n", HttpPort)
	err = server.Start()
	if err != nil {
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
