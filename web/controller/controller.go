package controller

import (
	"net/http"

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
	gatewayList, err := model.GetGateway()
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
	nodeList, err := model.GetTopology(ctx.Param("gateway"))
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
	if len(ctx.Param("advanced")) > 0 && ctx.Param("advanced") == "1" {
		nwStatDataAdv, err := model.GetNWStatAdv(ctx.Param("gateway"))
		if err != nil {
			return ctx.JSON(500, 0, err.Error(), nil)
		}
		if len(nwStatDataAdv) == 0 {
			return ctx.JSON(200, 0, "no result found", nil)
		}
		return ctx.JSON(200, 1, "success", nwStatDataAdv)
	}

	nwStatData, err := model.GetNWStat(ctx.Param("gateway"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nwStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nwStatData)
}

// GetSensorNWStat handles GET /api/:gateway/nwstat/:sensorID
func GetSensorNWStat(ctx *sweetygo.Context) error {
	sensorNWStatData, err := model.GetSensorNWStat(ctx.Param("gateway"), ctx.Param("sensorID"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(sensorNWStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", sensorNWStatData)
}
