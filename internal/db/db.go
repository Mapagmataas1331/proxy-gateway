package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB     *sql.DB
	Mutex  sync.Mutex
	dbFile = "/data/logs.db"
)

func Init() {
	os.MkdirAll(filepath.Dir(dbFile), 0755)

	var err error
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS ips (id INTEGER PRIMARY KEY, ip TEXT UNIQUE);
	CREATE TABLE IF NOT EXISTS paths (id INTEGER PRIMARY KEY, path TEXT UNIQUE);
	CREATE TABLE IF NOT EXISTS actions (id INTEGER PRIMARY KEY, action TEXT UNIQUE);
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time INTEGER NOT NULL,
		ip_id INTEGER NOT NULL,
		path_id INTEGER NOT NULL,
		port TEXT NOT NULL,
		protocol TEXT NOT NULL,
		auth_ok BOOLEAN NOT NULL,
		action_id INTEGER NOT NULL,
		FOREIGN KEY(ip_id) REFERENCES ips(id),
		FOREIGN KEY(path_id) REFERENCES paths(id),
		FOREIGN KEY(action_id) REFERENCES actions(id)
	);
	CREATE INDEX IF NOT EXISTS idx_time ON logs(time DESC);
	CREATE INDEX IF NOT EXISTS idx_auth_ok ON logs(auth_ok);
	CREATE INDEX IF NOT EXISTS idx_ip_id ON logs(ip_id);
	CREATE INDEX IF NOT EXISTS idx_path_id ON logs(path_id);
	`

	if _, err := DB.Exec(schema); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}
}
