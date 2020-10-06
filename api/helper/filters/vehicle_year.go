package filters

import (
	"errors"
	"maga-auctions/entity"
)

type vehicleYear struct {
	ModelYear         int
	ManufacturingYear int
}

// NewVehicleYear filters vehicles by year of manufacture and model year
func NewVehicleYear(model, manufacturing int) (Filter, error) {
	if manufacturing > model {
		return nil, errors.New("year of manufacture cannot be greater than the model")
	}

	return &vehicleYear{
		ManufacturingYear: manufacturing,
		ModelYear:         model,
	}, nil
}

// Rule filter model
func (v vehicleYear) Rule(vehicle entity.Vehicle) bool {
	return v.ManufacturingYear == vehicle.ManufacturingYear && v.ModelYear == vehicle.ModelYear
}

// Apply filter
func (v vehicleYear) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
