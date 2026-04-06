package repo

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	if err := createMigrationsTable(db); err != nil {
		return err
	}

	for _, file := range migrationFiles {
		if err := runMigration(db, file); err != nil {
			return err
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func getMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"
	entries, err := fs.ReadDir(os.DirFS(migrationsDir), ".")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, filepath.Join(migrationsDir, entry.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

func runMigration(db *sql.DB, filename string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE filename = ?", filename).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Migration %s already applied", filename)
		return nil
	}

	log.Printf("Running migration: %s", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO schema_migrations (filename) VALUES (?)", filename)
	return err
}
