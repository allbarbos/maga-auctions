package filters_test

import (
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVehicleYear_Rule(t *testing.T) {
	t.Run("must validate vehicle by model letters", func(t *testing.T) {
		ft, _ := filters.NewVehicleYear(2000, 1999)
		ve := entity.Vehicle{ManufacturingYear: 1999, ModelYear: 2000}

		assert.True(t, ft.Rule(ve))
	})
}

func TestNewVehicleYear_RuleError(t *testing.T) {
	t.Run("must validate vehicle by model letters", func(t *testing.T) {
		_, err := filters.NewVehicleYear(1999, 2000)

		assert.EqualError(t, err, "year of manufacture cannot be greater than the model")
	})
}

func TestNewVehicleYear_Apply(t *testing.T) {
	t.Run("must filter by model letters", func(t *testing.T) {
		ve1 := entity.Vehicle{ManufacturingYear: 1999, ModelYear: 2000}
		ve2 := entity.Vehicle{ManufacturingYear: 2019, ModelYear: 2020}
		items := &[]entity.Vehicle{ve1, ve2}

		ft, _ := filters.NewVehicleYear(2000, 1999)
		ft.Apply(items)

		assert.Len(t, *items, 1)
	})
}
