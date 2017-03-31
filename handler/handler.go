package handler

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/faststackco/machinestack/scheduler"
	"github.com/go-pg/pg"
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
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	AppMetadata AppMetadata `json:"app_metadata"`
	jwt.StandardClaims
}

// AppMetadata do we really need this? TODO
type AppMetadata struct {
	Quota Quota `json:"quota"`
}

// Quota defines how many instances a user can create, and how much cpu and ram they get
// CPU translates to cores, RAM translates to GB
type Quota struct {
	Instances int `json:"instances"`
	CPU       int `json:"cpus"`
	RAM       int `json:"ram"`
}

// Machine is a machine
type Machine struct {
	Name   string
	Image  string // Image name
	Owner  string // Owner name
	Driver string // Driver name
	Node   string // Node ID
}
