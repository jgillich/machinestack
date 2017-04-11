package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

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
