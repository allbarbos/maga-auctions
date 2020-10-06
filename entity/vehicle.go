package entity

// Vehicle entity
type Vehicle struct {
	ID                int    `json:"id"`                // ID - Identificador único do veículo
	Brand             string `json:"brand"`             // MARCA - Marca do veículo
	Model             string `json:"model"`             // MODELO - Modelo do veículo
	ModelYear         int    `json:"modelYear"`         // ANOMODELO - Ano do modelo do veículo
	ManufacturingYear int    `json:"manufacturingYear"` // ANOFABRICACAO - Ano de fabricação do veículo
	Lot               Lot    `json:"lot"`
	Bid               Bid    `json:"bid"`
}

type VehiclesAsc []Vehicle

func (v VehiclesAsc) Len() int           { return len(v) }
func (v VehiclesAsc) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v VehiclesAsc) Less(i, j int) bool { return v[i].Bid.Date.Before(v[j].Bid.Date) }

type VehiclesDesc []Vehicle

func (v VehiclesDesc) Len() int           { return len(v) }
func (v VehiclesDesc) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v VehiclesDesc) Less(i, j int) bool { return v[i].Bid.Date.After(v[j].Bid.Date) }
