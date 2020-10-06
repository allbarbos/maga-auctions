package controller_test

import (
	"maga-auctions/api/controller"
	"maga-auctions/legacy"
	mock_legacy "maga-auctions/legacy/mocks"
	"maga-auctions/utils"
	"maga-auctions/vehicle"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func mockApiLegacy(pathJSON string, statusCode int) {
	legacy.APIURI = "https://test.com"
	legacy.Client = &mock_legacy.MockClient{}
	mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{Body: utils.TestMakeBody(pathJSON), StatusCode: statusCode}, nil
	}
}

func TestHealthCheck(t *testing.T) {
	t.Run("must return status ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewHealthCheck(srv).HealthCheck(c)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"status":"ok","dependencies":{"legacyApi":"ok"}}`, w.Body.String())
	})
}

func TestHealthCheck_Errors(t *testing.T) {
	t.Run("must return status error when any dependency fail", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		mockApiLegacy("", 500)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewHealthCheck(srv).HealthCheck(c)

		assert.Equal(t, 500, w.Code)
		assert.JSONEq(t, `{"status":"error","dependencies":{"legacyApi":"error"}}`, w.Body.String())
	})
}
