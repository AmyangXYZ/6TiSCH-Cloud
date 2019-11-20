package controller

import (
	"net/http"
	"time"

	"../model"
	"github.com/AmyangXYZ/sweetygo"
)

// Index page handler.
func Index(ctx *sweetygo.Context) error {
	return ctx.Render(200, "index")
}

// Static files handler.
func Static(ctx *sweetygo.Context) error {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
	return nil
}

// GetGateway handles GET /api/gateway
func GetGateway(ctx *sweetygo.Context) error {
	timeRange := range2stamp(ctx.Param("range"))
	gatewayList, err := model.GetGateway(timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(gatewayList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", gatewayList)
}

// GetTopology handles GET /api/:gateway/topology
func GetTopology(ctx *sweetygo.Context) error {
	timeRange := range2stamp(ctx.Param("range"))
	nodeList, err := model.GetTopology(ctx.Param("gateway"), timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nodeList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nodeList)
}

// GetNWStat handles GET /api/:gateway/nwstat
func GetNWStat(ctx *sweetygo.Context) error {
	timeRange := range2stamp(ctx.Param("range"))
	nwStatData, err := model.GetNWStat(ctx.Param("gateway"), timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nwStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nwStatData)
}

// GetNWStatByID handles GET /api/:gateway/nwstat/:sensorID
func GetNWStatByID(ctx *sweetygo.Context) error {
	timeRange := range2stamp(ctx.Param("range"))
	if ctx.Param("advanced") != "" && ctx.Param("advanced") == "1" {
		sensorNWStatAdvData, err := model.GetNWStatAdvByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange)
		if err != nil {
			return ctx.JSON(500, 0, err.Error(), nil)
		}
		if len(sensorNWStatAdvData) == 0 {
			return ctx.JSON(200, 0, "no result found", nil)
		}
		return ctx.JSON(200, 1, "success", sensorNWStatAdvData)
	}

	sensorNWStatData, err := model.GetNWStatByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(sensorNWStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", sensorNWStatData)
}

// GetBattery handles GET /api/:gateway/battery
func GetBattery(ctx *sweetygo.Context) error {
	timeRange := range2stamp(ctx.Param("range"))
	batData, err := model.GetBattery(ctx.Param("gateway"), timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(batData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", batData)
}

func range2stamp(timeRange string) int64 {
	now := time.Now().UnixNano() / 1e6
	startTime := int64(0)
	switch timeRange {
	case "hour":
		startTime = now - 60*60*1000
	case "day":
		startTime = now - 60*60*24*1000
	case "week":
		startTime = now - 60*60*24*7*1000
	case "month":
		startTime = now - 60*60*24*7*30*1000
	default:
		break
	}
	return startTime
}
