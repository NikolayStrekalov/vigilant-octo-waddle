package pgstorage

import (
	"context"
	"fmt"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
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
	ret, err := p.db.GetGauges(context.TODO())
	if err != nil {
		logger.Info("error while query gauges:", err)
	}
	return ret
}

func (p *PGStorage) GetCounterList() []CounterListItem {
	ret, err := p.db.GetCounters(context.TODO())
	if err != nil {
		logger.Info("error while query counters:", err)
	}
	return ret
}

func (p *PGStorage) GetGauge(name string) (float64, error) {
	return p.db.GetGauge(context.TODO(), name)
}

func (p *PGStorage) GetCounter(name string) (int64, error) {
	return p.db.GetCounter(context.TODO(), name)
}

func (p *PGStorage) UpdateGauge(name string, value float64) {
	err := p.db.UpdateGauge(context.TODO(), name, value)
	if err != nil {
		logger.Info("failed to update gauge:", err)
	}
}

func (p *PGStorage) IncrementCounter(name string, value int64) {
	err := p.db.IncrementCounter(context.TODO(), name, value)
	if err != nil {
		logger.Info("failed to update counter:", err)
	}
}
func (p *PGStorage) BulkUpdate(metrics models.MetricsSlice) {
	err := p.db.BulkUpdate(context.TODO(), metrics)
	if err != nil {
		logger.Info("failed doing bulk update:", err)
	}
}
