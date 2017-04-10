package model

import (
	"time"

	"github.com/faststackco/machinestack/config"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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

	for _, model := range []interface{}{&Machine{}} {
		if err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true, // TODO remove temp, use migrations
		}); err != nil {
			return nil, err
		}
	}

	return db, nil
}
