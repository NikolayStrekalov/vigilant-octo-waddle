package pgstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN string
}

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	pool, err := initPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	db := &DB{
		pool: pool,
	}
	err = db.initTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init tables: %w", err)
	}
	return db, nil
}

func initPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return pool, nil
}

func (db *DB) initTables(ctx context.Context) error {
	row := db.pool.QueryRow(ctx, "SELECT version FROM migrations WHERE id = 1")
	var version int
	err := row.Scan(&version)
	if err != nil {
		return fmt.Errorf("error reading schema version: %w", err)
	}
	if version == 0 {
		err := db.createSchema(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}

func (db *DB) createSchema(ctx context.Context) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				logger.Info("failed to rollback the transaction", err)
			}
		}
	}()

	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS migrations(
			id INT PRIMARY KEY,
			version INT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS gauges(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name VARCHAR(200) UNIQUE NOT NULL,
			value DOUBLE PRECISION NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS counters(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name VARCHAR(200) UNIQUE NOT NULL,
			value BIGINT NOT NULL
		)`,
		`INSERT INTO migrations(id, version) VALUES (1, 1)`,
	}

	for _, stmt := range createSchemaStmts {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return nil
}

func (db *DB) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64
	row := db.pool.QueryRow(ctx, "SELECT value FROM gauges WHERE name=$1;", name)
	err := row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting gauge '%s': %w", name, err)
	}
	return value, nil
}

func (db *DB) GetGauges(ctx context.Context) ([]GaugeListItem, error) {
	var (
		name  string
		value float64
		ret   = []GaugeListItem{}
	)
	rows, err := db.pool.Query(ctx, "SELECT name, value FROM gauges;")
	if err != nil {
		return ret, fmt.Errorf("error fetching gauges: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			return ret, fmt.Errorf("error reading gauges: %w", err)
		}
		ret = append(ret, GaugeListItem{Name: name, Value: value})
	}
	return ret, nil
}

func (db *DB) GetCounter(ctx context.Context, name string) (int64, error) {
	var value int64
	row := db.pool.QueryRow(ctx, "SELECT value FROM counters WHERE name=$1;", name)
	err := row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error getting counter '%s': %w", name, err)
	}
	return value, nil
}

func (db *DB) GetCounters(ctx context.Context) ([]CounterListItem, error) {
	var (
		name  string
		value int64
		ret   = []CounterListItem{}
	)
	rows, err := db.pool.Query(ctx, "SELECT name, value FROM counters;")
	if err != nil {
		return ret, fmt.Errorf("error fetching counters: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			return ret, fmt.Errorf("error reading counters: %w", err)
		}
		ret = append(ret, CounterListItem{Name: name, Value: value})
	}
	return ret, nil
}

func (db *DB) UpdateGauge(ctx context.Context, name string, value float64) error {
	sql := `INSERT INTO gauges(name, value) VALUES ($1, $2)
	ON CONFLICT ON CONSTRAINT gauges_name_key DO UPDATE SET value = EXCLUDED.value;`
	_, err := db.pool.Exec(ctx, sql, name, value)
	if err != nil {
		return fmt.Errorf("failed to update gauge %s: %w", name, err)
	}
	return nil
}

func (db *DB) IncrementCounter(ctx context.Context, name string, value int64) error {
	sql := `INSERT INTO counters(name, value) VALUES ($1, $2)
	ON CONFLICT ON CONSTRAINT counters_name_key DO UPDATE SET value = counters.value + EXCLUDED.value;`
	_, err := db.pool.Exec(ctx, sql, name, value)
	if err != nil {
		return fmt.Errorf("failed to increment counter %s: %w", name, err)
	}
	return nil
}
