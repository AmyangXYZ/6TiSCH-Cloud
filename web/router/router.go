package router

import (
	"github.com/AmyangXYZ/6TiSCH-Cloud/web/controller"
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter .
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", controller.Index)
	app.GET("/static/*files", controller.Static)

	app.GET("/api/gateway", controller.GetGateway)
	app.GET("/api/:gateway/topology", controller.GetTopology)
	app.GET("/api/:gateway/topology/history", controller.GetTopoHistory)
	app.GET("/api/:gateway/schedule", controller.GetSchedule)
	app.GET("/api/:gateway/schedule/partition", controller.GetPartition)
	app.GET("/api/:gateway/schedule/partition_harp", controller.GetPartition)
	app.GET("/api/:gateway/nwstat", controller.GetNWStat)
	app.GET("/api/:gateway/nwstat/:sensorID", controller.GetNWStatByID)
	app.GET("/api/:gateway/nwstat/:sensorID/latency", controller.GetLatencyByID)
	app.GET("/api/:gateway/nwstat/:sensorID/channel", controller.GetChInfoByID)
	app.GET("/api/:gateway/battery", controller.GetBattery)
	app.GET("/api/:gateway/battery/:sensorID", controller.GetBatteryByID)
	app.GET("/api/:gateway/noise", controller.GetNoiseLevel)
	app.GET("/api/:gateway/txtotal", controller.GetTxTotal)
}
