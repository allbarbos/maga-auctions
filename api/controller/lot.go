package controller

import (
	"context"
	"maga-auctions/api/handler"
	"maga-auctions/vehicle"
	"time"

	"github.com/gin-gonic/gin"
)

// LotController contract
type LotController interface {
	VehiclesByLot(c *gin.Context)
}

type lotCtrl struct {
	srv vehicle.Service
}

// NewLot controller
func NewLot(srv vehicle.Service) LotController {
	return &lotCtrl{
		srv: srv,
	}
}

func (v lotCtrl) VehiclesByLot(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	vs, err := v.srv.ByLotID(ctx, id, c.Query("bidOrder"))

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	handler.ResponseSuccess(200, vs, c)
}
