package api

import (
	"testing"

	"gitlab.com/faststack/machinestack/scheduler"

	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
)

var (
	testToken = jwt.Token{Claims: jwt.MapClaims{
		"id":            1,
		"machine_quota": map[string]int{"instances": 10, "cpu": 1, "ram": 1},
	}}
	testScheduler, _ = scheduler.NewMockScheduler(nil)
	mockScheduler    = testScheduler.(*scheduler.MockScheduler)
	testDB           *pg.DB
	testHandler      *Handler
)

func TestMain(m *testing.M) {
	testDB = pg.Connect(&pg.Options{
		Addr:     os.Getenv("POSTGRES_ADDR"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	})

	if _, err := testDB.Query(nil, "TRUNCATE TABLE machines;"); err != nil {
		panic(err)
	}

	testHandler = &Handler{
		DB:           testDB,
		Scheduler:    testScheduler,
		AllowOrigins: []string{},
		JWTSecret:    []byte("foobar"),
	}

	retCode := m.Run()

	testDB.Close()
	os.Exit(retCode)
}
