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

type vehicleResponse struct {
	ID             int     `json:"ID"`
	DataLance      string  `json:"DATALANCE"`
	Lote           string  `json:"LOTE"`
	CodigoControle string  `json:"CODIGOCONTROLE"`
	Marca          string  `json:"MARCA"`
	Modelo         string  `json:"MODELO"`
	AnoFabricacao  int     `json:"ANOFABRICACAO"`
	AnoModelo      int     `json:"ANOMODELO"`
	ValorLance     float32 `json:"VALORLANCE"`
	UsuarioLance   string  `json:"USUARIOLANCE"`
}

func validResponseBody(payload interface{}) io.ReadCloser {
	b, _ := json.Marshal(payload)
	return ioutil.NopCloser(bytes.NewReader(b))
}

func invalidResponseBody() io.ReadCloser {
	items := struct{ invalid bool }{invalid: true}

	b, _ := json.Marshal(items)
	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestGet(t *testing.T) {
	t.Run("must return a list of vehicles", func(t *testing.T) {
		items := []vehicleResponse{
			{
				ID:             1,
				DataLance:      "21/08/2020 - 11:24",
				Lote:           "0033",
				CodigoControle: "80623",
				Marca:          "YAMAHA",
				Modelo:         "T115 CRYPTON ED",
				AnoFabricacao:  2011,
				AnoModelo:      2011,
				ValorLance:     0,
				UsuarioLance:   "-",
			},
		}

		legacy.Client = &mock_legacy.MockClient{}
		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{Body: validResponseBody(items)}, nil
		}

		api := legacy.NewAPI()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

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
			legacy.APIURI = tt.apiURI
			legacy.Client = &mock_legacy.MockClient{}
			mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
				return tt.doRes, tt.doError
			}

			api := legacy.NewAPI()
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			_, err := api.Get(ctx)

			assert.NotNil(t, err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("must register a vehicle", func(t *testing.T) {
		legacy.Client = &mock_legacy.MockClient{}
		v := vehicleResponse{
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
		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{Body: validResponseBody(v)}, nil
		}

		api := legacy.NewAPI()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		a := entity.Vehicle{
			Brand:             "RENAULT",
			Model:             "CLIO 16VS",
			ModelYear:         2007,
			ManufacturingYear: 2007,
			Lot: entity.Lot{
				ID:           "0196",
				VehicleLotID: "56248",
			},
		}

		resp, _ := api.Create(ctx, a)

		assert.NotNil(t, resp)
		assert.EqualValues(t, v, *resp)
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
			a := entity.Vehicle{
				Brand:             "RENAULT",
				Model:             "CLIO 16VS",
				ModelYear:         2007,
				ManufacturingYear: 2007,
				Lot: entity.Lot{
					ID:           "0196",
					VehicleLotID: "56248",
				},
			}

			_, err := api.Create(ctx, a)

			assert.NotNil(t, err)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

// func TestUpdate(t *testing.T) {
// 	t.Run("must update a vehicle", func(t *testing.T) {
// 		legacy.Client = &mock_legacy.MockClient{}
// 		mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
// 			return &http.Response{Body: http.NoBody, StatusCode: 200}, nil
// 		}

// 		api := legacy.NewAPI()
// 		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 		defer cancel()

// 		a := entity.Auction{
// 			LotID:        "0196",
// 			VehicleLotID: "56248",
// 			Bid: entity.Bid{
// 				Date:  "21/08/2020 - 13:24",
// 				User:  "-",
// 				Value: 0,
// 			},
// 			Vehicle: entity.Vehicle{
// 				ID:                1,
// 				Brand:             "RENAULT",
// 				ManufacturingYear: 2007,
// 				Model:             "CLIO 16VS",
// 				ModelYear:         2007,
// 			},
// 		}

// 		resp, _ := api.Update(ctx, 1, a)

// 		assert.NotNil(t, resp)
// 		assert.Equal(t, 200, resp.StatusCode)
// 	})
// }

// func TestUpdate_Errors(t *testing.T) {
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
// 			a := entity.Auction{
// 				LotID:        "0196",
// 				VehicleLotID: "56248",
// 				Bid: entity.Bid{
// 					Date:  "21/08/2020 - 13:24",
// 					User:  "-",
// 					Value: 0,
// 				},
// 				Vehicle: entity.Vehicle{
// 					ID:                1,
// 					Brand:             "RENAULT",
// 					ManufacturingYear: 2007,
// 					Model:             "CLIO 16VS",
// 					ModelYear:         2007,
// 				},
// 			}

// 			_, err := api.Update(ctx, 1, a)

// 			assert.NotNil(t, err)
// 			assert.Equal(t, tt.want, err.Error())
// 		})
// 	}
// }

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
