package entity

// Lot entity
type Lot struct {
	ID           string `json:"id"`           // LOTE - Agrupador de um conjunto de veículos
	VehicleLotID string `json:"vehicleLotId"` // CODIGOCONTROLE - Código único do veículo dentro do lote
}
