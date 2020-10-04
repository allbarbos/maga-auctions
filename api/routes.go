package api

import (
	ctrl "maga-auctions/api/controller"
	"maga-auctions/legacy"
	"maga-auctions/utils"
	"maga-auctions/vehicle"

	"github.com/gin-gonic/gin"
)

// Config routes
func Config() *gin.Engine {
	if utils.EnvVars.API.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/maga-auctions/v1/health-check", healthCtrl().HealthCheck)
	router.POST("/maga-auctions/v1/vehicles", vehicleCtrl().Create)
	router.GET("/maga-auctions/v1/vehicles/:id", vehicleCtrl().ByID)
	router.PUT("/maga-auctions/v1/vehicles/:id", vehicleCtrl().Update)

	return router
}

func healthCtrl() ctrl.HealthCheck {
	return ctrl.HealthCheck{}
}

func vehicleCtrl() ctrl.VehicleController {
	api := legacy.NewAPI()
	srv := vehicle.NewService(api)
	return ctrl.NewVehicle(srv)
}
