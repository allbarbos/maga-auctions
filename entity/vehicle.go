package entity

// Vehicle entity
type Vehicle struct {
	ID                int    `json:"id"`                // ID - Identificador único do veículo
	Brand             string `json:"brand"`             // MARCA - Marca do veículo
	Model             string `json:"model"`             // MODELO - Modelo do veículo
	ModelYear         int    `json:"modelYear"`         // ANOMODELO - Ano do modelo do veículo
	ManufacturingYear int    `json:"manufacturingYear"` // ANOFABRICACAO - Ano de fabricação do veículo
	Lot               Lot    `json:"lot"`
}
