package handler

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/faststackco/machinestack/scheduler"
	"github.com/go-pg/pg"
)

type Handler struct {
	db    *pg.DB
	sched scheduler.Scheduler
}

func NewHandler(db *pg.DB, sched scheduler.Scheduler) *Handler {
	return &Handler{db, sched}
}

type JwtClaims struct {
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	AppMetadata AppMetadata `json:"app_metadata"`
	jwt.StandardClaims
}

type AppMetadata struct {
	Quota Quota `json:"quota"`
}

type Quota struct {
	Instances int `json:"instances"`
	Cpus      int `json:"cpus"`
	Ram       int `json:"ram"`
}

type Machine struct {
	Name   string
	Image  string // Image name
	Owner  string // Owner name
	Driver string // Driver name
	Node   string // Node ID
}
