package filters_test

import (
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVehicleModel_Rule(t *testing.T) {
	t.Run("must validate vehicle by model letters", func(t *testing.T) {
		ft := filters.NewVehicleModel("CLI")
		ve := entity.Vehicle{Model: "CLIO 16VS"}

		assert.True(t, ft.Rule(ve))
	})
}

func TestNewVehicleModel_Apply(t *testing.T) {
	t.Run("must filter by model letters", func(t *testing.T) {
		ve1 := entity.Vehicle{Model: "CLIO 16VS"}
		ve2 := entity.Vehicle{Model: "VERSA 16SL CVT"}
		items := &[]entity.Vehicle{ve1, ve2}

		ft := filters.NewVehicleModel("CLI")
		ft.Apply(items)

		assert.Len(t, *items, 1)
	})
}
