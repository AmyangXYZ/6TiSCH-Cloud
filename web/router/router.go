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
	app.GET("/api/:gatewayName/topology", controller.GetTopologyData)
	app.GET("/api/:gatewayName/nwstat", controller.GetNetworkStat)
}
