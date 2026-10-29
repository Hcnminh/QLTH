package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	globalDB  *sql.DB
	globalErr error
	initOnce  sync.Once
	dbLock    sync.RWMutex
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// GetDB returns the PostgreSQL database connection pool.
func GetDB() (*sql.DB, error) {
	dbLock.RLock()
	if globalDB != nil {
		defer dbLock.RUnlock()
		return globalDB, nil
	}
	dbLock.RUnlock()

	initOnce.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			host := getEnv("PGHOST", "localhost")
			port := getEnv("PGPORT", "5432")
			user := getEnv("PGUSER", "qlth_user")
			pwd := getEnv("PGPASSWORD", "qlth_password_secure")
			dbname := getEnv("PGDATABASE", "qlth_eduprimary")
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pwd, host, port, dbname)
		}

		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			globalErr = err
			return
		}

		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(10)
		conn.SetConnMaxLifetime(5 * time.Minute)

		if err := conn.Ping(); err != nil {
			globalErr = err
			return
		}

		dbLock.Lock()
		globalDB = conn
		dbLock.Unlock()
	})

	return globalDB, globalErr
}

// SetDB allows injecting a custom or mock database connection (useful in testing)
func SetDB(d *sql.DB) {
	dbLock.Lock()
	defer dbLock.Unlock()
	globalDB = d
	globalErr = nil
}
