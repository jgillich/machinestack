package api

import (
	"net/http"

	"gitlab.com/faststack/machinestack/scheduler"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/google/jsonapi"
	"github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

var (
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
	// UnauthorizedError is returned when request is not authenticated
	UnauthorizedError = &jsonapi.ErrorObject{
		Code:   "unauthorized",
		Title:  "Unauthorized",
		Detail: "The request lacks valid authentication credentials.",
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
	// TODO configurable
	logger, _ = zap.NewDevelopment()
)

// Handler stores common types needed by the api
type Handler struct {
	DB           *pg.DB
	Scheduler    scheduler.Scheduler
	JWTSecret    []byte
	AllowOrigins []string
}

// WriteOneError returns one error object
func WriteOneError(w http.ResponseWriter, status int, err *jsonapi.ErrorObject) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	if err := jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{err}); err != nil {
		logger.Error("error marshalling json response", zap.Error(err), zap.Any("err", err))
	}
}

// WriteOne returns one resource object
func WriteOne(w http.ResponseWriter, status int, model interface{}) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	if err := jsonapi.MarshalOnePayload(w, model); err != nil {
		logger.Error("error marshalling json response", zap.Error(err), zap.Any("model", model))
	}
}

// WriteMany returns a list of resource objects
func WriteMany(w http.ResponseWriter, status int, models interface{}) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	if err := jsonapi.MarshalManyPayload(w, models); err != nil {
		logger.Error("error marshalling json response", zap.Error(err), zap.Any("models", models))
	}
}

// WriteInternalError logs and returns an internal server error
func WriteInternalError(w http.ResponseWriter, log string, err error) {
	logger.Error(log, zap.Error(err))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusInternalServerError)
	if err := jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{InternalServerError}); err != nil {
		logger.Error("error marshalling json response", zap.Error(err))
	}
}

// Serve sets up the http server and listens on the provided addr
func (h *Handler) Serve(addr string) error {

	router := httprouter.New()
	router.GET("/machines", h.MachineList)
	router.POST("/machines", h.MachineCreate)
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
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	jwt := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return h.JWTSecret, nil
		},
		SigningMethod:       jwt.SigningMethodHS256,
		CredentialsOptional: true,
	})
	middleware.Use(negroni.HandlerFunc(jwt.HandlerWithNext))

	middleware.UseHandler(router)

	return http.ListenAndServe(addr, middleware)
}
