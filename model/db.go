package model

import (
	"time"

	"github.com/faststackco/machinestack/config"
	"github.com/go-pg/pg"
)

// Db creates a database connection and applies migrations
func Db(config *config.PostgresConfig) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:        config.Address,
		User:        config.Username,
		Password:    config.Password,
		Database:    config.Database,
		PoolSize:    20,
		PoolTimeout: time.Second * 5,
		ReadTimeout: time.Second * 5,
	})

	return db, nil
}
