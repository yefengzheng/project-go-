package pgsql // Consider renaming this to 'mysql' for clarity

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"project-go-/internal/config"
	"time"
)

type Context struct {
	DB *sql.DB
}

// CreateNewPgsqlContext initializes and connects to the MySQL database
func CreateNewPgsqlContext(cfg *config.Config, lifeTime time.Duration) (*Context, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&parseTime=true",
		cfg.PGSQL.User,
		cfg.PGSQL.Password,
		cfg.PGSQL.Address,
		cfg.PGSQL.Port,
		cfg.PGSQL.ResultDb,
		cfg.PGSQL.ConnectTimeout,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Set connection pool options
	db.SetConnMaxLifetime(lifeTime)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Ping to test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &Context{DB: db}, nil
}

// Optional: Retrieve image data by name
func (ctx *Context) GetImage(name string) ([]byte, error) {
	fmt.Printf("[DEBUG] Looking up: '%s'\n", name) // Add this
	var data []byte
	err := ctx.DB.QueryRow("SELECT data FROM images WHERE name = ?", name).Scan(&data)
	if err != nil {
		fmt.Printf("[DEBUG] Not found: %v\n", err)
	}
	return data, err
}
