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
	code := 200

	if err != nil {
		s = "error"
		code = 500
	}

	hc := entity.HealthCheck{
		Status: s,
		Dependencies: map[string]string{
			"legacyApi": s,
		},
	}

	c.JSON(code, hc)
}
