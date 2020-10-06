package filters

import (
	"errors"
	"maga-auctions/entity"
)

// VehicleYearBetween filters vehicles by year of manufacture
type vehicleYearBetween struct {
	Min int
	Max int
}

// NewVehicleYearBetween filter
func NewVehicleYearBetween(min, max int) (Filter, error) {
	if max < min {
		return nil, errors.New("year of manufacture max cannot be less than min")
	}

	return &vehicleYearBetween{
		Min: min,
		Max: max,
	}, nil
}

// Rule filter model
func (v vehicleYearBetween) Rule(vehicle entity.Vehicle) bool {
	return v.Min <= vehicle.ManufacturingYear && v.Max >= vehicle.ManufacturingYear
}

// Apply filter
func (v vehicleYearBetween) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
