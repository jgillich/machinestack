package model

import (
	"fmt"

	"github.com/go-pg/migrations"
)

// Machine is a machine, obviously
type Machine struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Driver string `json:"driver"`
	Owner  string `json:"-"`
	Node   string `json:"-"`
}

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table machines...")
		_, err := db.Exec(`
			CREATE TABLE machines (
				id bigserial,
				name text,
				image text,
				driver text,
				owner text,
				node text,
				PRIMARY KEY (id)
			)
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table machines...")
		_, err := db.Exec(`DROP TABLE machines`)
		return err
	})
}
