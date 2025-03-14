package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ivolejon/pivo/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func config() *pgxpool.Config {
	env := settings.Environment()

	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(env.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

type DB struct {
	Pool *pgxpool.Pool
}

var (
	DbInstance *DB
	Once       sync.Once
)

func ConnectAndGetPool(ctx context.Context) (*DB, error) {
	var err error

	Once.Do(func() {
		db, dbErr := pgxpool.NewWithConfig(context.Background(), config())
		if dbErr != nil {
			err = fmt.Errorf("unable to create connection pool: %w", dbErr)
			return
		}

		DbInstance = &DB{db}
	})
	return DbInstance, err
}

func (db *DB) Ping(ctx context.Context) error {
	p, _ := db.Pool.Acquire(ctx)
	return p.Ping(ctx)
}

func (db *DB) Close() {
	db.Pool.Close()
}
