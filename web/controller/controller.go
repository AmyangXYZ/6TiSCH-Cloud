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
// return a gateway list
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

// GetTopologyData handles GET /api/topology
// return node
func GetTopologyData(ctx *sweetygo.Context) error {
	nodeList, err := model.GetTopologyData(ctx.Param("gatewayName"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nodeList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nodeList)
}

// GetNetworkStat handles GET /api/:gatewayName/nwstat
// return an array of sensor_id and avg_rtt.
// if advanced=1, return [sensor_id, avg_mac_tx_total_diff,
// avg_mac_tx_noack_diff,avg_app_per_sent_diff, avg_app_per_lost_diff]
func GetNetworkStat(ctx *sweetygo.Context) error {
	if len(ctx.Param("advanced")) > 0 && ctx.Param("advanced") == "1" {
		nwStatDataAdv, err := model.GetNWStatDataAdv(ctx.Param("gatewayName"))
		if err != nil {
			return ctx.JSON(500, 0, err.Error(), nil)
		}
		if len(nwStatDataAdv) == 0 {
			return ctx.JSON(200, 0, "no result found", nil)
		}
		return ctx.JSON(200, 1, "success", nwStatDataAdv)
	}

	nwStatData, err := model.GetNWStatData(ctx.Param("gatewayName"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nwStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}

	return ctx.JSON(200, 1, "success", nwStatData)
}
