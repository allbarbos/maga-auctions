package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ResponseSuccess creates payload
func ResponseSuccess(status int, body interface{}, c *gin.Context) {
	if body != nil {

		c.JSON(status, body)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

// ResponseError creates payload
func ResponseError(err error, c *gin.Context) {
	typeError := reflect.TypeOf(err).String()
	status := http.StatusInternalServerError
	message := err.Error()

	switch typeError {
	case "handler.BadRequest":
		status = http.StatusBadRequest
	case "handler.NotFound":
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{"error": message})
}
