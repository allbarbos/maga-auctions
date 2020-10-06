package controller

import (
	"context"
	"maga-auctions/entity"
	"maga-auctions/vehicle"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck contract
type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

type healthCheck struct {
	srv vehicle.Service
}

// NewHealthCheck controller
func NewHealthCheck(srv vehicle.Service) HealthCheck {
	return &healthCheck{
		srv: srv,
	}
}

// HealthCheck returns application health
func (h healthCheck) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := h.srv.ByID(ctx, 1)

	s := "ok"

	if err != nil {
		s = "error"
	}

	hc := entity.HealthCheck{
		Status: "ok",
		Dependencies: map[string]string{
			"legacyApi": s,
		},
	}

	c.JSON(200, hc)
}
