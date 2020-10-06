package entity

import "time"

// Bid entity
type Bid struct {
	Date  time.Time `json:"date"`  // DATALANCE - Data/hora que foi realizado o último lance
	Value float32   `json:"value"` // VALORLANCE - Valor do último lance dado
	User  string    `json:"user"`  // USUARIOLANCE - Usuário cadastrado na plataforma que fez o último lance
}
