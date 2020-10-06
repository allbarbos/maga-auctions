package controller_test

import (
	"bytes"
	"maga-auctions/api/controller"
	"maga-auctions/legacy"
	"maga-auctions/vehicle"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("must create a new vehicle", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := bytes.NewBufferString(`{"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}}`)
		c.Request, _ = http.NewRequest("POST", "/vehicles", body)

		mockApiLegacy("testdata/criar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewVehicle(srv).Create(c)

		assert.Equal(t, 201, w.Code)
		assert.Equal(t, w.HeaderMap["Location"][0], "/9999")
		assert.JSONEq(
			t,
			`{"vehicle":{"id":9999,"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}},"links":[{"uri":"/9999","rel":"self","type":"GET"},{"uri":"/9999","rel":"self","type":"PUT"},{"uri":"/9999","rel":"self","type":"DELETE"}]}`,
			w.Body.String(),
		)
	})
}

func TestCreate_Errors(t *testing.T) {
	testCases := []struct {
		desc, body, jsonPATH, wantJson string
		wantStatus                     int
	}{
		{
			desc:       "must return error when body is invalid",
			body:       "{",
			wantStatus: 400,
			wantJson:   `{"error":"body is invalid"}`,
		},
		{
			desc:       "must return error when legacy api fails",
			body:       `{"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}}`,
			jsonPATH:   "testdata/criar_response_error_api.json",
			wantStatus: 500,
			wantJson:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewBufferString(tt.body)
			c.Request, _ = http.NewRequest("POST", "/vehicles", body)
			mockApiLegacy(tt.jsonPATH, 200)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewVehicle(srv).Create(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(
				t,
				tt.wantJson,
				w.Body.String(),
			)
		})
	}
}

func TestAll(t *testing.T) {
	t.Run("must return a list of vehicle", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/vehicles?bidOrder=desc&brand=renault&model=S&manufacturingYearMin=2011&manufacturingYearMax=2013&manufacturingYear=2011&modelYear=2011", nil)
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewVehicle(srv).All(c)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(
			t,
			`[{"id":544,"brand":"RENAULT","model":"SYMBOL EX1616V","modelYear":2011,"manufacturingYear":2011,"lot":{"id":"0046","vehicleLotId":"716797"},"bid":{"date":"2020-08-22T09:48:00Z","value":4500,"user":"Sorico1"}}]`,
			w.Body.String(),
		)
	})
}

func TestAll_Errors(t *testing.T) {
	testCases := []struct {
		desc, jsonPATH, query, wantJson string
		wantStatus                      int
	}{
		{
			desc:       "must return error when manufacturing year min is invalid",
			query:      "/vehicles?manufacturingYearMin=a",
			wantStatus: 400,
			wantJson:   `{"error":"manufacturing year min is invalid"}`,
		},
		{
			desc:       "must return error when manufacturing year max is invalid",
			query:      "/vehicles?manufacturingYearMin=2011&manufacturingYearMax=b",
			wantStatus: 400,
			wantJson:   `{"error":"manufacturing year max is invalid"}`,
		},
		{
			desc:       "must return error when year of manufacture max cannot be less than min",
			query:      "/vehicles?manufacturingYearMin=2011&manufacturingYearMax=1999",
			wantStatus: 400,
			wantJson:   `{"error":"year of manufacture max cannot be less than min"}`,
		},
		{
			desc:       "must return error when manufacturing year is invalid",
			query:      "/vehicles?manufacturingYearMin=2011&manufacturingYearMax=2011&manufacturingYear=a",
			wantStatus: 400,
			wantJson:   `{"error":"manufacturing year is invalid"}`,
		},
		{
			desc:       "must return error when model year is invalid",
			query:      "/vehicles?manufacturingYearMin=2011&manufacturingYearMax=2011&manufacturingYear=2011&modelYear=a",
			wantStatus: 400,
			wantJson:   `{"error":"model year is invalid"}`,
		},
		{
			desc:       "must return error when model year is invalid",
			query:      "/vehicles?manufacturingYearMin=2011&manufacturingYearMax=2011&manufacturingYear=3000&modelYear=2012",
			wantStatus: 400,
			wantJson:   `{"error":"year of manufacture cannot be greater than the model"}`,
		},
		{
			desc:       "must return error when legacy api fails",
			jsonPATH:   "testdata/consultar_response_error_api.json",
			wantStatus: 500,
			wantJson:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", tt.query, nil)
			mockApiLegacy(tt.jsonPATH, 200)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewVehicle(srv).All(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantJson, w.Body.String())
		})
	}
}

func TestByID(t *testing.T) {
	t.Run("must return vehicle", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "760"}}
		c.Request, _ = http.NewRequest("GET", "http://test.com", nil)
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewVehicle(srv).ByID(c)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(
			t,
			`{"vehicle":{"id":760,"brand":"IVECO","model":"EUROCARGO 260E25N","modelYear":2012,"manufacturingYear":2011,"lot":{"id":"0068","vehicleLotId":"126845"},"bid":{"date":"2020-08-27T10:20:00Z","value":75000,"user":"ALDOBARROSO"}},"links":[{"uri":"","rel":"self","type":"PUT"},{"uri":"","rel":"self","type":"DELETE"}]}`,
			w.Body.String(),
		)
	})
}

func TestByID_Errors(t *testing.T) {
	testCases := []struct {
		desc, jsonPATH, id, wantJson string
		wantStatus                   int
	}{
		{
			desc:       "must return error when id is invalid",
			id:         "a",
			wantStatus: 400,
			wantJson:   `{"error":"id is invalid"}`,
		},
		{
			desc:       "must return error when legacy api fails",
			id:         "760",
			jsonPATH:   "testdata/consultar_response_error_api.json",
			wantStatus: 500,
			wantJson:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			c.Request, _ = http.NewRequest("GET", "http://test.com", nil)
			mockApiLegacy(tt.jsonPATH, 200)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewVehicle(srv).ByID(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantJson, w.Body.String())
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Run("must update a new vehicle", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "760"}}
		body := bytes.NewBufferString(`{"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}}`)
		c.Request, _ = http.NewRequest("PUT", "/vehicles", body)

		mockApiLegacy("testdata/alterar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewVehicle(srv).Update(c)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(
			t,
			`{"vehicle":{"id":760,"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}},"links":[{"uri":"","rel":"self","type":"GET"},{"uri":"","rel":"self","type":"DELETE"}]}`,
			w.Body.String(),
		)
	})
}

func TestUpdate_Errors(t *testing.T) {
	testCases := []struct {
		desc, id, body, wantJson string
		wantStatus               int
	}{
		{
			desc:       "must return error when id is invalid",
			id:         "a",
			wantStatus: 400,
			wantJson:   `{"error":"id is invalid"}`,
		},
		{
			desc:       "must return error when body is invalid",
			id:         "760",
			body:       "{",
			wantStatus: 400,
			wantJson:   `{"error":"body is invalid"}`,
		},
		{
			desc:       "must return error when legacy api fails",
			id:         "760",
			body:       `{"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}}`,
			wantStatus: 400,
			wantJson:   `{"error":"id not found"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			body := bytes.NewBufferString(tt.body)
			c.Request, _ = http.NewRequest("PUT", "/vehicles", body)
			mockApiLegacy("testdata/alterar_response_error_api.json", 200)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewVehicle(srv).Update(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantJson, w.Body.String())
		})
	}
}

func TestDelete(t *testing.T) {
	t.Run("must delete vehicle", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "760"}}

		mockApiLegacy("testdata/apagar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewVehicle(srv).Delete(c)

		assert.Equal(t, 200, w.Code)
	})
}

func TestDelete_Errors(t *testing.T) {
	testCases := []struct {
		desc, id, body, wantJson string
		wantStatus               int
	}{
		{
			desc:       "must return error when id is invalid",
			id:         "a",
			wantStatus: 400,
			wantJson:   `{"error":"id is invalid"}`,
		},
		{
			desc:       "must return error when legacy api fails",
			id:         "760",
			body:       `{"brand":"RENAULT","model":"CLIO 16VS","modelYear":2007,"manufacturingYear":2007,"lot":{"id":"0196","vehicleLotId":"56248"},"bid":{"date":"2020-08-27T10:20:00Z","value":15000,"user":"ALLBARBOS"}}`,
			wantStatus: 400,
			wantJson:   `{"error":"id not found"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			mockApiLegacy("testdata/apagar_response_error_api.json", 200)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewVehicle(srv).Delete(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantJson, w.Body.String())
		})
	}
}
