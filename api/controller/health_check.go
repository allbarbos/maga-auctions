package controller

import (
	"maga-auctions/entity"

	"github.com/gin-gonic/gin"
)

// HealthCheck controller
type HealthCheck struct{}

// HealthCheck returns application health
func (h HealthCheck) HealthCheck(c *gin.Context) {
	hc := entity.HealthCheck{
		Status: "ok",
	}

	c.JSON(200, hc)
}
