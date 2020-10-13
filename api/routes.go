package api

import (
	"fmt"
	ctrl "maga-auctions/api/controller"
	"maga-auctions/api/middlewares"
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

	app := gin.New()
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.Use(middlewares.CORS())
	app.NoRoute(middlewares.NoRouteHandler())

	app.GET("/maga-auctions/v1/health-check", healthCtrl().HealthCheck)

	app.POST("/maga-auctions/v1/vehicles", vehicleCtrl().Create)
	app.GET("/maga-auctions/v1/vehicles", vehicleCtrl().All)
	app.GET("/maga-auctions/v1/vehicles/:id", vehicleCtrl().ByID)
	app.PUT("/maga-auctions/v1/vehicles/:id", vehicleCtrl().Update)
	app.DELETE("/maga-auctions/v1/vehicles/:id", vehicleCtrl().Delete)

	app.GET("/maga-auctions/v1/lots/:id/vehicles", lotCtrl().VehiclesByLot)

	return app
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
