package legacy_test

import (
	"maga-auctions/entity"
	"maga-auctions/legacy"
	mock_legacy "maga-auctions/legacy/mocks"
	"maga-auctions/utils"

	"context"
	"errors"
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
)

func mockApiLegacy(apiURI, pathJSON string, statusCode int, doError error) {
	legacy.APIURI = apiURI
	legacy.Client = &mock_legacy.MockClient{}
	mock_legacy.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{Body: utils.TestMakeBody(pathJSON), StatusCode: statusCode}, doError
	}
}

func TestGet(t *testing.T) {
	t.Run("must return a list of vehicles", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("https://test.com", "testdata/consultar_response_api.json", 200, nil)

		api := legacy.NewAPI()
		resp, err := api.Get(ctx)

		assert.Nil(t, err)
		assert.Equal(t, len(resp), 764)
	})
}

func TestGet_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, jsonPATH, want string
		legacyApiStatusCode          int
		doError                      error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttps://test.com",
			want:   `parse "h!ttps://test.com": first path segment in URL cannot contain colon`,
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			doError:             errors.New("an error occurred while requesting the legacy api"),
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when the object parse fails",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 200,
			jsonPATH:            "testdata/consultar_response_error_api.json",
			want:                "invalid character '}' looking for beginning of value",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.apiURI, tt.jsonPATH, tt.legacyApiStatusCode, tt.doError)

			api := legacy.NewAPI()
			_, err := api.Get(ctx)

			assert.NotNil(t, err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("must register a vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("https://test.com", "testdata/criar_response_api.json", 200, nil)

		api := legacy.NewAPI()
		err := api.Create(ctx, &ve)

		assert.Nil(t, err)
	})
}

func TestCreate_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, jsonPATH, want string
		legacyApiStatusCode          int
		doError                      error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttps://test.com",
			want:   `parse "h!ttps://test.com": first path segment in URL cannot contain colon`,
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			doError:             errors.New("an error occurred while requesting the legacy api"),
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when the object parse fails",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 200,
			jsonPATH:            "testdata/criar_response_error_api.json",
			want:                "invalid character '}' looking for beginning of value",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.apiURI, tt.jsonPATH, tt.legacyApiStatusCode, tt.doError)

			api := legacy.NewAPI()
			err := api.Create(ctx, &ve)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Run("must update a vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("https://test.com", "testdata/alterar_response_api.json", 200, nil)

		api := legacy.NewAPI()
		err := api.Update(ctx, &ve)

		assert.Nil(t, err)
	})
}

func TestUpdate_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, jsonPATH, want string
		legacyApiStatusCode          int
		doError                      error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttp://error",
			want:   `parse "h!ttp://error": first path segment in URL cannot contain colon`,
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			doError:             errors.New("an error occurred while requesting the legacy api"),
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when failing the legacy api request",
			apiURI:              "https://test.com",
			legacyApiStatusCode: 500,
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when id is not found",
			legacyApiStatusCode: 200,
			jsonPATH:            "testdata/apagar_response_error_api.json",
			want:                "id not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			mockApiLegacy(tt.apiURI, tt.jsonPATH, tt.legacyApiStatusCode, tt.doError)

			api := legacy.NewAPI()
			err := api.Update(ctx, &ve)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Run("must delete vehicle", func(t *testing.T) {
		defer cancel()
		mockApiLegacy("https://test.com", "testdata/apagar_response_api.json", 200, nil)

		api := legacy.NewAPI()
		err := api.Delete(ctx, 1)

		assert.Nil(t, err)
	})
}

func TestDelete_Errors(t *testing.T) {
	testCases := []struct {
		desc, apiURI, jsonPATH, want string
		legacyApiStatusCode          int
		doError                      error
	}{
		{
			desc:   "must return error for an invalid url",
			apiURI: "h!ttp://error",
			want:   `parse "h!ttp://error": first path segment in URL cannot contain colon`,
		},
		{
			desc:                "must return an error when failing the legacy api request",
			legacyApiStatusCode: 500,
			doError:             errors.New("an error occurred while requesting the legacy api"),
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when failing the legacy api request",
			legacyApiStatusCode: 500,
			want:                "an error occurred while requesting the legacy api",
		},
		{
			desc:                "must return an error when id is not found",
			jsonPATH:            "testdata/apagar_response_error_api.json",
			legacyApiStatusCode: 200,
			want:                "id not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			defer cancel()
			mockApiLegacy(tt.apiURI, tt.jsonPATH, tt.legacyApiStatusCode, tt.doError)

			api := legacy.NewAPI()
			err := api.Delete(ctx, 1)

			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.want)
		})
	}
}
