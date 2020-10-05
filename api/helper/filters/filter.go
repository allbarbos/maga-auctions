package filters

import (
	"maga-auctions/entity"
)

// Filter interface
type Filter interface {
	Apply(input *[]entity.Vehicle)
	Rule(vehicle entity.Vehicle) bool
}

func filterApply(input *[]entity.Vehicle, rule func(entity.Vehicle) bool) {
	items := *input
	output := []entity.Vehicle{}

	for i := 0; i < len(items); i++ {
		vehicle := items[i]

		if rule(vehicle) {
			output = append(output, vehicle)
		}
	}

	*input = output
}
