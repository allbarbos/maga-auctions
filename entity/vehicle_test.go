package entity_test

import (
	"maga-auctions/entity"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	v1 = entity.Vehicle{
		ID:                1,
		Brand:             "brand",
		Model:             "model",
		ManufacturingYear: 1000,
		ModelYear:         1000,
		Bid: entity.Bid{
			Date:  time.Now(),
			User:  "test",
			Value: 4.3,
		},
		Lot: entity.Lot{
			ID:           "id",
			VehicleLotID: "id",
		},
	}

	v2 = entity.Vehicle{
		ID:                2,
		Brand:             "brand",
		Model:             "model",
		ManufacturingYear: 1000,
		ModelYear:         1000,
		Bid: entity.Bid{
			Date:  time.Now().Add(5 * time.Hour),
			User:  "test",
			Value: 4.3,
		},
		Lot: entity.Lot{
			ID:           "id",
			VehicleLotID: "id",
		},
	}
)

func TestVehicleAsc(t *testing.T) {
	items := []entity.Vehicle{v2, v1}

	assert.Equal(t, items[0].ID, 2)

	sort.Sort(entity.VehiclesAsc(items))

	assert.Equal(t, items[0].ID, 1)
}

func TestVehicleDesc(t *testing.T) {
	items := []entity.Vehicle{v1, v2}

	assert.Equal(t, items[0].ID, 1)

	sort.Sort(entity.VehiclesDesc(items))

	assert.Equal(t, items[0].ID, 2)
}
