package vehicle

import (
	"maga-auctions/api/handler"
	"maga-auctions/entity"
	"maga-auctions/legacy"

	"context"
)

// Service contract
type Service interface {
	Create(ctx context.Context, vehicle entity.Vehicle) (*entity.Vehicle, error)
	ByID(ctx context.Context, id int) (*entity.Vehicle, error)
	Update(ctx context.Context, vehicle *entity.Vehicle) error
	Delete(ctx context.Context, id int) error
}

type srv struct {
	legacyAPI legacy.API
}

// NewService returns a planet service instance
func NewService(api legacy.API) Service {
	return &srv{
		legacyAPI: api,
	}
}

func (s srv) ByID(ctx context.Context, id int) (*entity.Vehicle, error) {
	if id <= 0 {
		return nil, handler.BadRequest{Message: "invalid id"}
	}

	items, err := s.legacyAPI.Get(ctx)

	if err != nil || items == nil {
		return nil, handler.InternalServer{Message: "error when searching for vehicles in legacy api"}
	}

	if id > len(items) {
		return nil, handler.BadRequest{Message: "invalid id"}
	}

	l := items[id-1]
	v := &entity.Vehicle{
		ID:                l.ID,
		Brand:             l.Marca,
		Model:             l.Modelo,
		ModelYear:         l.AnoModelo,
		ManufacturingYear: l.AnoFabricacao,
		Lot: entity.Lot{
			ID:           l.Lote,
			VehicleLotID: l.CodigoControle,
		},
	}

	return v, nil
}

func (s srv) Create(ctx context.Context, vehicle entity.Vehicle) (*entity.Vehicle, error) {
	err := s.legacyAPI.Create(ctx, &vehicle)

	if err != nil || vehicle.ID == 0 {
		return nil, handler.InternalServer{Message: err.Error()}
	}

	return &vehicle, nil
}

func (s srv) Update(ctx context.Context, vehicle *entity.Vehicle) error {
	if vehicle.ID <= 0 {
		return handler.BadRequest{Message: "invalid id"}
	}

	if err := s.legacyAPI.Update(ctx, vehicle); err != nil {
		msg := err.Error()
		if msg == "id not found" {
			return handler.BadRequest{Message: msg}
		}

		return handler.InternalServer{Message: msg}
	}

	return nil
}

func (s srv) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return handler.BadRequest{Message: "invalid id"}
	}

	if err := s.legacyAPI.Delete(ctx, id); err != nil {
		msg := err.Error()
		if msg == "id not found" {
			return handler.BadRequest{Message: msg}
		}

		return handler.InternalServer{Message: msg}
	}

	return nil
}
