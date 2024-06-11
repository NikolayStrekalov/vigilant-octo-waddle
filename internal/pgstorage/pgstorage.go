package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGStorage struct {
	db *DB
}

type GaugeListItem = struct {
	Name  string
	Value float64
}

type CounterListItem = struct {
	Name  string
	Value int64
}

const maxRequestAttempts = 4

var RetryOptions = []retry.Option{
	retry.RetryIf(func(err error) bool {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return true
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgerrcode.IsConnectionException(pgErr.Code)
		}
		return false
	}),
	retry.Attempts(maxRequestAttempts),
	retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		return time.Duration(1+n*2) * time.Second
	}),
}

func NewPGStorage(dsn string) (*PGStorage, func() error, error) {
	db, err := NewDB(context.TODO(), Config{DSN: dsn})
	if err != nil {
		return nil, func() error { return db.Close() }, fmt.Errorf("init db error: %w", err)
	}
	return &PGStorage{
		db: db,
	}, func() error { return db.Close() }, nil
}

func (p *PGStorage) GetGaugeList() []GaugeListItem {
	ret, err := retry.DoWithData(
		func() ([]GaugeListItem, error) {
			return p.db.GetGauges(context.TODO())
		},
		RetryOptions...,
	)
	if err != nil {
		logger.Info("error while query gauges:", err)
	}
	return ret
}

func (p *PGStorage) GetCounterList() []CounterListItem {
	ret, err := retry.DoWithData(
		func() ([]CounterListItem, error) {
			return p.db.GetCounters(context.TODO())
		},
		RetryOptions...,
	)
	if err != nil {
		logger.Info("error while query counters:", err)
	}
	return ret
}

func (p *PGStorage) GetGauge(name string) (float64, error) {
	val, err := retry.DoWithData(
		func() (float64, error) {
			return p.db.GetGauge(context.TODO(), name)
		},
		RetryOptions...,
	)
	if err != nil {
		return val, fmt.Errorf("failed to get gauge %s: %w", name, err)
	}
	return val, nil
}

func (p *PGStorage) GetCounter(name string) (int64, error) {
	val, err := retry.DoWithData(
		func() (int64, error) {
			return p.db.GetCounter(context.TODO(), name)
		},
		RetryOptions...,
	)
	if err != nil {
		return val, fmt.Errorf("failed to get counter %s: %w", name, err)
	}
	return val, nil
}

func (p *PGStorage) UpdateGauge(name string, value float64) {
	err := retry.Do(
		func() error {
			return p.db.UpdateGauge(context.TODO(), name, value)
		},
		RetryOptions...,
	)
	if err != nil {
		logger.Info("failed to update gauge:", err)
	}
}

func (p *PGStorage) IncrementCounter(name string, value int64) {
	err := retry.Do(
		func() error {
			return p.db.IncrementCounter(context.TODO(), name, value)
		},
		RetryOptions...,
	)
	if err != nil {
		logger.Info("failed to update counter:", err)
	}
}

func (p *PGStorage) BulkUpdate(metrics models.MetricsSlice) {
	err := retry.Do(
		func() error {
			return p.db.BulkUpdate(context.TODO(), metrics)
		},
		RetryOptions...,
	)
	if err != nil {
		logger.Info("failed doing bulk update:", err)
	}
}
