package main

import "github.com/AmyangXYZ/sgo"

// SetRouter .
func SetRouter(app *sgo.SGo) {
	app.GET("/", Index)
	app.GET("/static/*files", Static)

	app.GET("/api/gateway", GetGateway)
	app.GET("/api/:gateway/topology", GetTopology)
	app.GET("/api/:gateway/topology/history", GetTopoHistory)
	app.GET("/api/:gateway/schedule", GetSchedule)
	app.GET("/api/:gateway/schedule/partition", GetPartition)
	app.GET("/api/:gateway/schedule/partition_harp", GetPartitionHARP)
	app.GET("/api/:gateway/schedule/sp_harp", GetSubPartitionHARP)
	app.GET("/api/:gateway/nwstat", GetNWStat)
	app.GET("/api/:gateway/nwstat/:sensorID", GetNWStatByID)
	app.GET("/api/:gateway/nwstat/:sensorID/latency", GetLatencyByID)
	app.GET("/api/:gateway/nwstat/:sensorID/channel", GetChInfoByID)
	app.GET("/api/:gateway/battery", GetBattery)
	app.GET("/api/:gateway/battery/:sensorID", GetBatteryByID)
	app.GET("/api/:gateway/noise", GetNoiseLevel)
	app.GET("/api/:gateway/txtotal", GetTxTotal)
	app.GET("/api/:gateway/:sensorID/sensors", GetSensorDataByID)
	app.PUT("/api/:gateway/:sensorID/sensors/:type", PutSensorRateByIDnType)
	app.OPTIONS("/api/:gateway/:sensorID/sensors/:type", sgo.PreflightHandler)
}
