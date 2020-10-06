package filters_test

import (
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVehicleBrand_Rule(t *testing.T) {
	t.Run("must validate the vehicle brand", func(t *testing.T) {
		ft := filters.NewVehicleBrand("brand")
		ve := entity.Vehicle{Brand: "BRAND"}

		assert.True(t, ft.Rule(ve))
	})
}

func TestVehicleBrand_Apply(t *testing.T) {
	t.Run("must filter by vehicle brand", func(t *testing.T) {
		ve1 := entity.Vehicle{Brand: "BRAND"}
		ve2 := entity.Vehicle{Brand: "BRAND2"}
		items := &[]entity.Vehicle{ve1, ve2}

		ft := filters.NewVehicleBrand("brand")
		ft.Apply(items)

		assert.Len(t, *items, 1)
	})
}
