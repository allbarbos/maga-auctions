package api

import (
	ctrl "maga-auctions/api/controller"
	"maga-auctions/legacy"
	"maga-auctions/utils"
	"maga-auctions/vehicle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config routes
func Config() *gin.Engine {
	if utils.EnvVars.API.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(configCors)

	router.GET("/maga-auctions/v1/health-check", healthCtrl().HealthCheck)

	router.POST("/maga-auctions/v1/vehicles", vehicleCtrl().Create)
	router.GET("/maga-auctions/v1/vehicles", vehicleCtrl().All)
	router.GET("/maga-auctions/v1/vehicles/:id", vehicleCtrl().ByID)
	router.PUT("/maga-auctions/v1/vehicles/:id", vehicleCtrl().Update)
	router.DELETE("/maga-auctions/v1/vehicles/:id", vehicleCtrl().Delete)

	router.GET("/maga-auctions/v1/lots/:id/vehicles", lotCtrl().VehiclesByLot)

	return router
}

func configCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func buildSrv() vehicle.Service {
	api := legacy.NewAPI()
	return vehicle.NewService(api)
}

func healthCtrl() ctrl.HealthCheck {
	return ctrl.NewHealthCheck(buildSrv())
}

func vehicleCtrl() ctrl.VehicleController {
	return ctrl.NewVehicle(buildSrv())
}

func lotCtrl() ctrl.LotController {
	return ctrl.NewLot(buildSrv())
}
