package vehicle_test

import (
	"context"
	"errors"
	"fmt"
	"maga-auctions/entity"
	"maga-auctions/legacy"
	mock_legacy "maga-auctions/legacy/mocks"
	"maga-auctions/vehicle"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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
)

func TestByID(t *testing.T) {
	t.Run("must return vehicle", func(t *testing.T) {
		defer cancel()
		c := gomock.NewController(t)
		defer c.Finish()

		api := mock_legacy.NewMockAPI(c)
		items := []legacy.VehicleLegacy{vehicleLegacy}

		api.EXPECT().Get(ctx).Return(items, nil)

		srv := vehicle.NewService(api)

		item, _ := srv.ByID(ctx, 1)

		assert.NotNil(t, item)
		assert.Equal(t, items[0].ID, item.ID)
		assert.Equal(t, items[0].Marca, item.Brand)
		assert.Equal(t, items[0].Modelo, item.Model)
		assert.Equal(t, items[0].AnoModelo, item.ModelYear)
		assert.Equal(t, items[0].AnoFabricacao, item.ManufacturingYear)
		assert.Equal(t, items[0].Lote, item.Lot.ID)
		assert.Equal(t, items[0].CodigoControle, item.Lot.VehicleLotID)
	})
}

func TestByID_Errors(t *testing.T) {
	testCases := []struct {
		desc, want string
		id         int
		err        error
		items      []legacy.VehicleLegacy
	}{
		{
			desc: "must return error when id less than zero",
			id:   0,
			want: "invalid id",
		},
		{
			desc: "must return error when legacy API fails",
			id:   1,
			err:  errors.New("legacy api error"),
			want: "internal server error",
		},
		{
			desc:  "must return error when id greater than zero",
			id:    2,
			want:  "invalid id",
			items: []legacy.VehicleLegacy{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			api := mock_legacy.NewMockAPI(c)

			if tt.items != nil || tt.err != nil {
				api.EXPECT().Get(ctx).Return(tt.items, tt.err)
			}

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
		c := gomock.NewController(t)
		defer c.Finish()
		v := entity.Vehicle{
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

		api := mock_legacy.NewMockAPI(c)
		api.EXPECT().Create(ctx, v).Return(&vehicleLegacy, nil)

		srv := vehicle.NewService(api)
		item, _ := srv.Create(ctx, v)

		assert.NotNil(t, item)
		assert.Equal(t, vehicleLegacy.ID, item.ID)
		assert.Equal(t, vehicleLegacy.Marca, item.Brand)
		assert.Equal(t, vehicleLegacy.Modelo, item.Model)
		assert.Equal(t, vehicleLegacy.AnoModelo, item.ModelYear)
		assert.Equal(t, vehicleLegacy.AnoFabricacao, item.ManufacturingYear)
		assert.Equal(t, vehicleLegacy.Lote, item.Lot.ID)
		assert.Equal(t, vehicleLegacy.CodigoControle, item.Lot.VehicleLotID)
	})
}

func TestCreate_Error(t *testing.T) {
	t.Run("must return error when legacy API fails", func(t *testing.T) {
		defer cancel()
		c := gomock.NewController(t)
		defer c.Finish()
		v := entity.Vehicle{
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

		api := mock_legacy.NewMockAPI(c)
		api.EXPECT().Create(ctx, v).Return(nil, errors.New("legacy API fails"))

		srv := vehicle.NewService(api)
		item, err := srv.Create(ctx, v)

		assert.Nil(t, item)
		assert.EqualError(t, err, "internal server error")
	})
}
