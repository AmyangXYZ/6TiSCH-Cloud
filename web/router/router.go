package router

import (
	"../controller"
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter .
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", controller.Index)
	app.GET("/static/*files", controller.Static)

	app.GET("/api/gateway", controller.GetGateway)
	app.GET("/api/:gateway/topology", controller.GetTopology)
	app.GET("/api/:gateway/nwstat", controller.GetNWStat)
	app.GET("/api/:gateway/nwstat/:sensorID", controller.GetNWStatByID)
	app.GET("/api/:gateway/battery", controller.GetBattery)
	app.GET("/api/:gateway/noise", controller.GetNoiseLevel)
}
