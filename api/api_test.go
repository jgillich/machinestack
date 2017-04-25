package api

import (
	"testing"

	"gitlab.com/faststack/machinestack/model"
	"gitlab.com/faststack/machinestack/scheduler"

	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	testToken = &jwt.Token{Claims: jwt.MapClaims{
		"id":            1,
		"machine_quota": map[string]interface{}{"instances": float64(10), "cpu": float64(1), "ram": float64(1)},
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

	if err := testDB.CreateTable(&model.Machine{}, &orm.CreateTableOptions{
		Temp: true,
	}); err != nil {
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
