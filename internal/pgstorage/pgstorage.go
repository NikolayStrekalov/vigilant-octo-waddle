package pgstorage

import (
	"context"
	"fmt"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
)

type PGStorage struct {
	db       *DB
	dumpFile string
	sync     bool
}

type GaugeListItem = struct {
	Name  string
	Value float64
}

type CounterListItem = struct {
	Name  string
	Value int64
}

func NewPGStorage(dsn string, synchronous bool, dumpPath string) (*PGStorage, func() error, error) {
	db, err := NewDB(context.TODO(), Config{DSN: dsn})
	if err != nil {
		return nil, func() error { return db.Close() }, fmt.Errorf("init db error: %w", err)
	}
	return &PGStorage{
		sync:     synchronous,
		dumpFile: dumpPath,
		db:       db,
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

func (p *PGStorage) Dump() {
	// Do nothing; dumped by DB
	// TODO: remove from StorageOperations?
}

func (p *PGStorage) Restore() {
	// Do nothing; restored by DB
	// TODO: remove from StorageOperations?
}