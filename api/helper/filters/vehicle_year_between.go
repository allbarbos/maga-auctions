package filters

import (
	"maga-auctions/entity"
)

// VehicleYearBetween filters vehicles by year of manufacture
type VehicleYearBetween struct {
	Min int
	Max int
}

// Rule filter model
func (v VehicleYearBetween) Rule(vehicle entity.Vehicle) bool {
	return v.Min <= vehicle.ManufacturingYear && v.Max >= vehicle.ManufacturingYear
}

// Apply filter
func (v VehicleYearBetween) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
