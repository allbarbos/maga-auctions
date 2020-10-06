package filters_test

import (
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVehicleYearBetween_Rule(t *testing.T) {
	t.Run("must validate the vehicle by year of manufacture", func(t *testing.T) {
		ft, _ := filters.NewVehicleYearBetween(2000, 2002)
		ve := entity.Vehicle{ManufacturingYear: 2001}

		assert.True(t, ft.Rule(ve))
	})
}

func TestNewVehicleYearBetween_RuleError(t *testing.T) {
	t.Run("must return error when year of manufacture max cannot be less than min", func(t *testing.T) {
		_, err := filters.NewVehicleYearBetween(2000, 1999)

		assert.EqualError(t, err, "year of manufacture max cannot be less than min")
	})
}

func TestNewVehicleYearBetween_Apply(t *testing.T) {
	t.Run("must filter by year of manufacture", func(t *testing.T) {
		ve1 := entity.Vehicle{ManufacturingYear: 2001}
		ve2 := entity.Vehicle{ManufacturingYear: 2019}
		items := &[]entity.Vehicle{ve1, ve2}

		ft, _ := filters.NewVehicleYearBetween(2000, 2002)
		ft.Apply(items)

		assert.Len(t, *items, 1)
	})
}
