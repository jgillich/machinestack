package handler

import (
	"fmt"

	"github.com/labstack/echo"
)

var (
	// StatusSuccess is the status of a successful response
	StatusSuccess = "success"
	// StatusError is the status of a error response
	StatusError = "error"
)

// Response is the structure of Handler responses
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Data responds with success and data
func Data(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, Response{
		Status: StatusSuccess,
		Data:   data,
	})
}

// Error responds with error and message
func Error(c echo.Context, status int, message string, a ...interface{}) error {
	return c.JSON(status, Response{
		Status:  StatusError,
		Message: fmt.Sprintf(message, a...),
	})
}

// Message responds with success and messsage
func Message(c echo.Context, status int, message string, a ...interface{}) error {
	return c.JSON(status, Response{
		Status:  StatusSuccess,
		Message: fmt.Sprintf(message, a...),
	})
}
