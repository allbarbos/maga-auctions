package filters

import "maga-auctions/entity"

// VehicleBrand filter
type VehicleBrand struct {
	Brand string
}

// Rule filter brand
func (v VehicleBrand) Rule(vehicle entity.Vehicle) bool {
	return vehicle.Brand == v.Brand
}

// Apply filter
func (v VehicleBrand) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
