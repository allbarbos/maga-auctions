package vehicle_test

import (
	"context"
	"errors"
	"fmt"
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"maga-auctions/legacy"
	mock_legacy "maga-auctions/legacy/mocks"
	"maga-auctions/utils"
	"maga-auctions/vehicle"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ctx, cancel   = context.WithTimeout(context.Background(), 2*time.Second)
	vehicleLegacy = legacy.VehicleLegacy{
		ID:             1,
		DataLance:      "21/08/2020 - 11:24",
		Lote:           "0033",
		CodigoControle: "80623",
		Marca:          "Marca Teste",
		Modelo:         "Modelo Teste",
		AnoFabricacao:  2011,
		AnoModelo:      2011,
		ValorLance:     0,
		UsuarioLance:   "-",
	}
	v = entity.Vehicle{
		ID:                vehicleLegacy.ID,
		Brand:             vehicleLegacy.Marca,
		Model:             vehicleLegacy.Modelo,
		ModelYear:         vehicleLegacy.AnoModelo,
		ManufacturingYear: vehicleLegacy.AnoFabricacao,
		Lot: entity.Lot{
			ID:           vehicleLegacy.Lote,
			VehicleLotID: vehicleLegacy.CodigoControle,
		},
	}
)

func mockApiLegacy(pathJSON string, statusCode int) {
	legacy.APIURI = "https://test.com"
	legacy.Client = &mock_legacy.MockClient{}
	mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{Body: utils.TestMakeBody(pathJSON), StatusCode: statusCode}, nil
	}
}

func TestAll(t *testing.T) {
	t.Run("must return list of vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)
		filters := []filters.Filter{
			filters.NewVehicleBrand("RENAULT"),
			filters.NewVehicleModel("S"),
		}

		resp, _ := srv.All(ctx, filters, "asc")

		assert.NotNil(t, resp)
		assert.Len(t, *resp, 20)
	})
}

func TestAll_Errors(t *testing.T) {
	t.Run("BILE", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("", 500)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		item, err := srv.All(ctx, []filters.Filter{}, "")

		assert.Nil(t, item)
		fmt.Print(err.Error())
		assert.EqualError(t, err, "internal server error")
	})
}

func TestByID(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("testdata/consultar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		item, _ := srv.ByID(ctx, 1)

		assert.NotNil(t, item)
		assert.Equal(t, 1, item.ID)
		assert.Equal(t, "RENAULT", item.Brand)
		assert.Equal(t, "CLIO 16VS", item.Model)
		assert.Equal(t, 2007, item.ModelYear)
		assert.Equal(t, 2007, item.ManufacturingYear)
		assert.Equal(t, "0196", item.Lot.ID)
		assert.Equal(t, "56248", item.Lot.VehicleLotID)
	})
}

func TestByID_Errors(t *testing.T) {
	testCases := []struct {
		desc, jsonPATH, want    string
		id, legacyApiStatusCode int
		err                     error
		items                   []legacy.VehicleLegacy
	}{
		{
			desc:                "must return error when id less than zero",
			id:                  0,
			legacyApiStatusCode: 500,
			want:                "invalid id",
		},
		{
			desc:                "must return error when legacy API fails",
			id:                  1,
			legacyApiStatusCode: 500,
			err:                 errors.New("legacy api error"),
			want:                "internal server error",
		},
		{
			desc:                "must return error when id greater than zero",
			id:                  9999,
			legacyApiStatusCode: 200,
			jsonPATH:            "testdata/consultar_response_api.json",
			want:                "invalid id",
			items:               []legacy.VehicleLegacy{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.jsonPATH, tt.legacyApiStatusCode)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			item, err := srv.ByID(ctx, tt.id)

			assert.Nil(t, item)
			fmt.Print(err.Error())
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("must create vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("testdata/criar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		item, _ := srv.Create(ctx, v)

		assert.NotNil(t, item)
		assert.Equal(t, 9999, item.ID)
		assert.Equal(t, v.Brand, item.Brand)
		assert.Equal(t, v.Model, item.Model)
		assert.Equal(t, v.ModelYear, item.ModelYear)
		assert.Equal(t, v.ManufacturingYear, item.ManufacturingYear)
		assert.Equal(t, v.Lot.ID, item.Lot.ID)
		assert.Equal(t, v.Lot.VehicleLotID, item.Lot.VehicleLotID)
	})
}

func TestCreate_Error(t *testing.T) {
	t.Run("must return error when legacy API fails", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("", 500)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		item, err := srv.Create(ctx, v)

		assert.Nil(t, item)
		assert.EqualError(t, err, "internal server error")
	})
}

func TestUpdate(t *testing.T) {
	t.Run("must update vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("testdata/alterar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		err := srv.Update(ctx, &v)

		assert.Nil(t, err)
	})
}

func TestUpdate_Errors(t *testing.T) {
	testCases := []struct {
		desc, jsonPATH, want    string
		id, legacyApiStatusCode int
		err                     error
		items                   []legacy.VehicleLegacy
	}{
		{
			desc: "must return error when id less than zero",
			id:   0,
			err:  errors.New("id not found"),
			want: "invalid id",
		},
		{
			desc:                "must return error when id is not found",
			id:                  1000,
			jsonPATH:            "testdata/alterar_response_error_api.json",
			legacyApiStatusCode: 200,
			err:                 errors.New("id not found"),
			want:                "id not found",
		},
		{
			desc:                "must return error when an unknown error occurs in legacy api",
			id:                  2,
			legacyApiStatusCode: 500,
			err:                 errors.New("error when updating in legacy api"),
			want:                "internal server error",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.jsonPATH, tt.legacyApiStatusCode)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			v.ID = tt.id
			err := srv.Update(ctx, &v)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Run("must delete vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("testdata/apagar_response_api.json", 200)

		api := legacy.NewAPI()
		srv := vehicle.NewService(api)

		err := srv.Delete(ctx, 1)

		assert.Nil(t, err)
	})
}

func TestDelete_Errors(t *testing.T) {
	testCases := []struct {
		desc, jsonPATH, want    string
		id, legacyApiStatusCode int
		err                     error
		items                   []legacy.VehicleLegacy
	}{
		{
			desc: "must return error when id less than zero",
			id:   0,
			err:  errors.New("id not found"),
			want: "invalid id",
		},
		{
			desc:                "must return error when id is not found",
			id:                  1000,
			legacyApiStatusCode: 200,
			jsonPATH:            "testdata/apagar_response_error_api.json",
			err:                 errors.New("id not found"),
			want:                "id not found",
		},
		{
			desc:                "must return error when an unknown error occurs in legacy api",
			id:                  2,
			legacyApiStatusCode: 500,
			err:                 errors.New("error when updating in legacy api"),
			want:                "internal server error",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.jsonPATH, tt.legacyApiStatusCode)

			api := legacy.NewAPI()
			srv := vehicle.NewService(api)

			v.ID = tt.id
			err := srv.Delete(ctx, tt.id)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}
