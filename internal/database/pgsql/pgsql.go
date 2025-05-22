package pgsql // Consider renaming this to 'mysql' for clarity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"project-go-/internal/config"
	"time"
)

type Context struct {
	DB *sql.DB
}

// CreateNewPgsqlContext initializes and connects to the MySQL database
func CreateNewPgsqlContext(cfg *config.Config, lifeTime time.Duration) (*Context, error) {
	log.Println("Connecting to PostgresSQL DB....")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PGSQL.User, cfg.PGSQL.Password, cfg.PGSQL.Address, cfg.PGSQL.Port, cfg.PGSQL.ResultDb)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Set connection pool options
	db.SetConnMaxLifetime(lifeTime)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Ping to test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres DB: %w", err)
	}

	log.Printf("Connected to PostgresDB %s", connStr)
	return &Context{DB: db}, nil
}

func (ctx *Context) InsertScanResult(req config.ScanRequest) error {
	// Convert []string to JSON for MySQL storage
	filesJSON, err := json.Marshal(req.MaliciousFiles)
	if err != nil {
		return fmt.Errorf("failed to marshal malicious files: %w", err)
	}

	query := `
		INSERT INTO scan_results (image_name, scan_start_time, scan_finish_time, scan_result, malicious_files)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = ctx.DB.Exec(query,
		req.ImageName,
		req.ScanStartTime,
		req.ScanFinishTime,
		req.ScanResult,
		string(filesJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert scan result: %w", err)
	}
	return nil
}

func (ctx *Context) GetScanResult(imageName string) (*config.ScanRequest, error) {
	query := `
		SELECT image_name, scan_start_time, scan_finish_time, scan_result, malicious_files
		FROM scan_results
		WHERE image_name = ?
		ORDER BY id DESC LIMIT 1
	`

	var scan config.ScanRequest
	var filesJSON string

	err := ctx.DB.QueryRow(query, imageName).Scan(
		&scan.ImageName,
		&scan.ScanStartTime,
		&scan.ScanFinishTime,
		&scan.ScanResult,
		&filesJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(filesJSON), &scan.MaliciousFiles); err != nil {
		return nil, err
	}

	return &scan, nil
}
