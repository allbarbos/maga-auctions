package filters

import (
	"maga-auctions/entity"
	"strings"
)

type vehicleBrand struct {
	Brand string
}

// NewVehicleBrand filter
func NewVehicleBrand(brand string) Filter {
	return &vehicleBrand{
		Brand: strings.ToUpper(brand),
	}
}

// Rule filter brand
func (v vehicleBrand) Rule(vehicle entity.Vehicle) bool {
	return vehicle.Brand == v.Brand
}

// Apply filter
func (v vehicleBrand) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
