package main

import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/s1ntezc0der/bazis-restapi/internal/config"
    "github.com/s1ntezc0der/bazis-restapi/pkg/db"
)

func migrateUp(db *sql.DB, dir string) error {
    files, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".sql") {
            continue
        }

        content, err := os.ReadFile(filepath.Join(dir, file.Name()))
        if err != nil {
            return err
        }

        if _, err := db.Exec(string(content)); err != nil {
            return fmt.Errorf("failed to apply %s: %w", file.Name(), err)
        }

        log.Printf("✅ Applied migration: %s", file.Name())
    }

    return nil
}

func migrateDown(db *sql.DB, dir string) error {
    log.Println("Down migration not fully implemented")
    return nil
}

func migrateCreate(dir, name string) error {
    timestamp := time.Now().Format("20060102150405")
    filename := filepath.Join(dir, timestamp+"_"+name+".sql")

    content := `
		-- +goose Up
		-- TODO: write your migration

		-- +goose Down
		-- TODO: write rollback
	`

    if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
        return err
    }

    log.Printf("Created migration: %s", filename)
    return nil
}

func main() {
	dir := flag.String("dir", "./migrations", "migrations directory")
	action := flag.String("action", "up", "up/down/create")
	name := flag.String("name", "", "migration name (for create)")
	flag.Parse()

	cfg := config.Load()

	conn, err := db.NewMySQLDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	switch *action {
	case "up":
		if err := migrateUp(conn, *dir); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := migrateDown(conn, *dir); err != nil {
			log.Fatal(err)
		}
	case "create":
		if *name == "" {
			log.Fatal("migration name is required")
		}
		if err := migrateCreate(*dir, *name); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown action")
	}
}

