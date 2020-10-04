package legacy_test

import (
	"maga-auctions/entity"
	"maga-auctions/legacy"
	mock_legacy "maga-auctions/legacy/mocks"

	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	ve          = entity.Vehicle{
		Brand:             "YAMAHA",
		Model:             "T115 CRYPTON ED",
		ModelYear:         2011,
		ManufacturingYear: 2011,
		Lot: entity.Lot{
			ID:           "0033",
			VehicleLotID: "80623",
		},
	}
	vel = legacy.VehicleLegacy{
		Lote:           "0033",
		CodigoControle: "80623",
		Marca:          "YAMAHA",
		Modelo:         "T115 CRYPTON ED",
		AnoFabricacao:  2011,
		AnoModelo:      2011,
		ID:             0,
		ValorLance:     0,
		DataLance:      "-",
		UsuarioLance:   "-",
	}
)

func validResponseBody(payload interface{}) io.ReadCloser {
	b, _ := json.Marshal(payload)
	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestGet(t *testing.T) {
	t.Run("must return a list of vehicles", func(t *testing.T) {
		defer cancel()
		items := []legacy.VehicleLegacy{vel}

		legacy.Client = &mock_legacy.MockClient{}
		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{Body: validResponseBody(items)}, nil
		}

		api := legacy.NewAPI()
		resp, err := api.Get(ctx)

		assert.Nil(t, err)
		assert.Greater(t, len(resp), 0)
	})
}

func TestGet_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI string
		doRes        *http.Response
		doError      error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttp://error",
		},
		{
			desc:    "must return an error when failing the legacy api request",
			doError: errors.New("error from legacy api"),
		},
		{
			desc:  "must return error when decode fails",
			doRes: &http.Response{Body: http.NoBody},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			legacy.APIURI = tt.apiURI
			legacy.Client = &mock_legacy.MockClient{}
			mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
				return tt.doRes, tt.doError
			}

			api := legacy.NewAPI()
			_, err := api.Get(ctx)

			assert.NotNil(t, err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("must register a vehicle", func(t *testing.T) {
		defer cancel()
		legacy.Client = &mock_legacy.MockClient{}
		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
			vel.ID = 13
			return &http.Response{Body: validResponseBody(vel)}, nil
		}

		api := legacy.NewAPI()
		err := api.Create(ctx, &ve)

		assert.Nil(t, err)
		assert.Greater(t, vel.ID, 0)
	})
}

func TestCreate_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, want string
		doRes              *http.Response
		doError            error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttp://error",
			want:   `parse "h!ttp://error": first path segment in URL cannot contain colon`,
		},
		{
			desc:    "must return an error when failing the legacy api request",
			doError: errors.New("error from legacy api"),
			want:    "error from legacy api",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			legacy.APIURI = tt.apiURI
			legacy.Client = &mock_legacy.MockClient{}
			mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
				return tt.doRes, tt.doError
			}

			api := legacy.NewAPI()
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			err := api.Create(ctx, &ve)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Run("must update a vehicle", func(t *testing.T) {
		defer cancel()
		legacy.Client = &mock_legacy.MockClient{}
		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{Body: validResponseBody(vel), StatusCode: 200}, nil
		}

		api := legacy.NewAPI()
		err := api.Update(ctx, &ve)

		assert.Nil(t, err)
	})
}

func TestUpdate_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, want string
		doRes              *http.Response
		doError            error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttp://error",
			want:   `parse "h!ttp://error": first path segment in URL cannot contain colon`,
		},
		{
			desc:    "must return an error when failing the legacy api request",
			doError: errors.New("error from legacy api"),
			want:    "error from legacy api",
		},
		{
			desc:  "must return an error when failing the legacy api request",
			doRes: &http.Response{StatusCode: 500},
			want:  "error when updating in legacy api",
		},
		{
			desc:  "must return an error when id is not found",
			doRes: &http.Response{StatusCode: 200, Body: validResponseBody(struct{ Mensagem string }{"nao encontrado"})},
			want:  "id not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			legacy.APIURI = tt.apiURI
			legacy.Client = &mock_legacy.MockClient{}
			mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
				return tt.doRes, tt.doError
			}

			api := legacy.NewAPI()
			err := api.Update(ctx, &ve)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}

// func TestDelete(t *testing.T) {
// 	t.Run("must delete vehicle", func(t *testing.T) {
// 		legacy.Client = &mock_legacy.MockClient{}
// 		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
// 			return &http.Response{Body: http.NoBody, StatusCode: 200}, nil
// 		}

// 		api := legacy.NewAPI()
// 		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 		defer cancel()

// 		resp, _ := api.Delete(ctx, 1)

// 		assert.NotNil(t, resp)
// 		assert.Equal(t, 200, resp.StatusCode)
// 	})
// }

// func TestDelete_Errors(t *testing.T) {
// 	testCases := []struct {
// 		desc, apiURI, want string
// 		doRes              *http.Response
// 		doError            error
// 	}{
// 		{
// 			desc:   "must return error for an invalid url",
// 			apiURI: "h!ttp://error",
// 			want:   `parse "h!ttp://error": first path segment in URL cannot contain colon`,
// 		},
// 		{
// 			desc:    "must return an error when failing the legacy api request",
// 			doError: errors.New("error from legacy api"),
// 			want:    "error from legacy api",
// 		},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			legacy.APIURI = tt.apiURI
// 			legacy.Client = &mock_legacy.MockClient{}
// 			mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
// 				return tt.doRes, tt.doError
// 			}

// 			api := legacy.NewAPI()
// 			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 			defer cancel()

// 			_, err := api.Delete(ctx, 1)

// 			assert.NotNil(t, err)
// 			assert.Equal(t, tt.want, err.Error())
// 		})
// 	}
// }
