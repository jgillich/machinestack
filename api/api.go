package api

import (
	"net/http"

	"gitlab.com/faststack/machinestack/scheduler"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/google/jsonapi"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

var (
	// UserContextKey is the context key for the jwt middleware
	UserContextKey = "user"
	// BadRequestError for badly formatted requests
	BadRequestError = &jsonapi.ErrorObject{
		Code:   "bad_request",
		Title:  "Bad Request",
		Detail: "Your request is not in a valid format",
	}
	// InternalServerError for unhandled errors
	InternalServerError = &jsonapi.ErrorObject{
		Code:   "internal_server",
		Title:  "Internal server error",
		Detail: "Unexpected internal error",
	}
	// ResourceNotFoundError for non existing resources
	ResourceNotFoundError = &jsonapi.ErrorObject{
		Code:   "not_found",
		Title:  "Resource was not found",
		Detail: "The resource you requested does not exist",
	}
	// AccessDeniedError is returned when rqeuested action is not allowed
	AccessDeniedError = &jsonapi.ErrorObject{
		Code:   "access_denied",
		Title:  "Access denied",
		Detail: "You are not allowed to perform the requested action.",
	}
	// ValidationFailedError is returned when request did not pass validation
	ValidationFailedError = &jsonapi.ErrorObject{
		Code:  "validation_failed",
		Title: "Request validation failed",
	}
	logger, _ = zap.NewProduction()
)

// Handler stores common types needed by the api
type Handler struct {
	DB           *pg.DB
	Scheduler    scheduler.Scheduler
	JWTSecret    []byte
	AllowOrigins []string
}

// JwtClaims are the custom claims we use
type JwtClaims struct {
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	MachineQuota MachineQuota `json:"machine_quota"`
	jwt.StandardClaims
}

// MachineQuota defines limits for a user
type MachineQuota struct {
	Instances int `json:"instances"`
	CPU       int `json:"cpu"`
	RAM       int `json:"ram"`
}

// WriteOneError returns one error object
func WriteOneError(w http.ResponseWriter, status int, err *jsonapi.ErrorObject) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	if err := jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{err}); err != nil {
		logger.Error("error marshalling json response", zap.Any("err", err))
	}
}

// WriteOne returns one resource object
func WriteOne(w http.ResponseWriter, status int, model interface{}) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	if err := jsonapi.MarshalOnePayload(w, model); err != nil {
		logger.Error("error marshalling json response", zap.Any("model", model))
	}
}

// WriteInternalError logs and returns an internal server error
func WriteInternalError(w http.ResponseWriter, log string, err error) {
	logger.Error(log, zap.Error(err))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusInternalServerError)
	if err := jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{InternalServerError}); err != nil {
		logger.Error("error marshalling json response", zap.Any("err", InternalServerError))
	}
}

// Serve sets up the http server and listens on the provided addr
func (h *Handler) Serve(addr string) error {

	router := httprouter.New()
	router.HandlerFunc("POST", "/machines", h.MachineCreate)
	router.DELETE("/machines/:name", h.MachineDelete)
	router.GET("/machines/:name", h.MachineInfo)
	router.POST("/machines/:name/session", h.SessionCreate)
	router.GET("/session/:id/io", h.SessionIO)
	router.GET("/session/:id/control", h.SessionControl)

	middleware := negroni.New()
	// TODO recovery middleware with jsonapi response
	middleware.Use(negroni.NewLogger())

	middleware.Use(cors.New(cors.Options{
		AllowedOrigins: h.AllowOrigins,
	}))

	jwt := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return h.JWTSecret, nil
		},
		SigningMethod:       jwt.SigningMethodHS256,
		CredentialsOptional: true,
		UserProperty:        UserContextKey,
	})
	middleware.Use(negroni.HandlerFunc(jwt.HandlerWithNext))

	middleware.UseHandler(router)

	return http.ListenAndServe(addr, middleware)
}
