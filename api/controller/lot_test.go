package controller_test

import (
	"maga-auctions/api/controller"
	"maga-auctions/legacy"
	"maga-auctions/vehicle"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVehiclesByLot(t *testing.T) {
	t.Run("must return status ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "0161"}}
		c.Request, _ = http.NewRequest("GET", "/vehicles?bidOrder=", nil)
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		controller.NewLot(srv).VehiclesByLot(c)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(
			t,
			`[{"id":180,"brand":"HONDA","model":"CIVIC SEDAN LXR","modelYear":2015,"manufacturingYear":2014,"lot":{"id":"0161","vehicleLotId":"733135"},"bid":{"date":"2020-08-21T12:58:00Z","value":5500,"user":"Michaelnf"}},{"id":725,"brand":"FIAT","model":"STRADA ADVENTURE CD","modelYear":2010,"manufacturingYear":2010,"lot":{"id":"0161","vehicleLotId":"733577"},"bid":{"date":"2020-08-22T11:15:00Z","value":22500,"user":"Damião A. d. S."}}]`,
			w.Body.String(),
		)
	})
}

func TestVehiclesByLot_Errors(t *testing.T) {
	testCases := []struct {
		desc, id, jsonPATH, wantJson    string
		legacyApiStatusCode, wantStatus int
	}{
		{
			desc:       "must return error when id is empty",
			id:         "",
			wantStatus: 400,
			wantJson:   `{"error":"invalid lot id"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			c.Request, _ = http.NewRequest("GET", "/vehicles?bidOrder=desc", nil)
			mockApiLegacy(tt.jsonPATH, tt.legacyApiStatusCode)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			controller.NewLot(srv).VehiclesByLot(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(
				t,
				tt.wantJson,
				w.Body.String(),
			)
		})
	}
}
