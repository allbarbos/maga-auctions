package filters

import (
	"maga-auctions/entity"
	"strings"
)

// VehicleModel filter
type VehicleModel struct {
	InitialLetters string
}

// Rule filter model
func (v VehicleModel) Rule(vehicle entity.Vehicle) bool {
	return strings.HasPrefix(vehicle.Model, v.InitialLetters)
}

// Apply filter
func (v VehicleModel) Apply(input *[]entity.Vehicle) {
	filterApply(input, v.Rule)
}
