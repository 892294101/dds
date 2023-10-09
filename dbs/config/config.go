package config

import "time"

type dbConfig struct {
	cacheFileSize    uint64
	cacheGroup       int
	flushPeriodTime  time.Duration
	flushMultipleOps int
}

type dbOption func(*dbConfig) error

func newConfig() *dbConfig {
	return new(dbConfig)
}

func SetCacheSize(cacheSize uint64) dbOption {
	return func(c *dbConfig) error {

		return nil
	}
}
