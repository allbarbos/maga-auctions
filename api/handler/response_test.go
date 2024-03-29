package handler_test

import (
	"maga-auctions/api/handler"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.ResponseSuccess(200, "ok", c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "\"ok\"", w.Body.String())
}

func TestResponseSuccess_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.ResponseSuccess(200, nil, c)

	assert.Equal(t, 200, w.Code)
}

func TestResponseError_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.ResponseError(handler.BadRequest{Message: "bad request"}, c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"bad request\"}", w.Body.String())
}

func TestResponseError_InternalServer(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.ResponseError(handler.InternalServer{Message: "internal server error"}, c)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"error\":\"internal server error\"}", w.Body.String())
}

func TestResponseError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.ResponseError(handler.NotFound{Message: "not found error"}, c)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"error\":\"not found error\"}", w.Body.String())
}
