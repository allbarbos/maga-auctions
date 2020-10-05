package filters

import (
	"maga-auctions/entity"
)

// VehicleYear filters vehicles by year of manufacture and model year
type VehicleYear struct {
	ModelYear         int
	ManufacturingYear int
}

// Rule filter model
func (v VehicleYear) Rule(vehicle entity.Vehicle) bool {
	return v.ManufacturingYear == vehicle.ManufacturingYear && v.ModelYear == vehicle.ModelYear
}

// Apply filter
func (v VehicleYear) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}

// // Validate manufacturing year and model year
// func (v VehicleYear) Validate() error {
// 	if v.ManufacturingYear <= 0 || v.ModelYear <= 0 {
// 		return errors.New("manufacturing year or model year invalid")
// 	}

// 	if v.ManufacturingYear > v.ModelYear {
// 		return errors.New("manufacturing year invalid")
// 	}

// 	return nil
// }
