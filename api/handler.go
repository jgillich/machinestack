package api

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/faststack/machinestack/scheduler"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

// Handler stores common types needed by the api
type Handler struct {
	db    *pg.DB
	sched scheduler.Scheduler
}

// NewHandler creates a new Handler
func NewHandler(db *pg.DB, sched scheduler.Scheduler) *Handler {
	return &Handler{db, sched}
}

// JwtClaims are the custom claims we use
type JwtClaims struct {
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	MachineQuota MachineQuota `json:"machine_quota"`
	jwt.StandardClaims
}

func getJwtClaims(c echo.Context) *JwtClaims {
	return c.Get("user").(*jwt.Token).Claims.(*JwtClaims)
}

// MachineQuota defines limits for a user
type MachineQuota struct {
	Instances int `json:"instances"`
	CPU       int `json:"cpu"`
	RAM       int `json:"ram"`
}
