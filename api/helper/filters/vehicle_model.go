package filters

import (
	"maga-auctions/entity"
	"strings"
)

type vehicleModel struct {
	InitialLetters string
}

// NewVehicleModel filter
func NewVehicleModel(letters string) Filter {
	return &vehicleModel{
		InitialLetters: strings.ToUpper(letters),
	}
}

// Rule filter model
func (v vehicleModel) Rule(vehicle entity.Vehicle) bool {
	return strings.HasPrefix(vehicle.Model, v.InitialLetters)
}

// Apply filter
func (v vehicleModel) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
