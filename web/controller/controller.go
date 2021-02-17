package controller

import (
	"net/http"
	"time"

	"github.com/AmyangXYZ/6TiSCH-Cloud/web/model"
	"github.com/AmyangXYZ/sweetygo"
)

var lastBootTime int64 = 0

func init() {
	go func() {
		for {
			lastBootTime = model.GetLastBootTime()
			time.Sleep(1 * time.Minute)
		}
	}()

}

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
	timeRange, now := range2stamp(ctx.Param("range"))
	gatewayList, err := model.GetGateway(timeRange, now)
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
	timeRange, now := range2stamp(ctx.Param("range"))
	nodeList, err := model.GetTopology(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nodeList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nodeList)
}

// GetTopoHistory handles GET /api/:gateway/topology/history
func GetTopoHistory(ctx *sweetygo.Context) error {
	eventList, err := model.GetTopoHistory(range2stamp("week"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(eventList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", eventList)
}

// GetSchedule handles GET /api/:gateway/schedule
func GetSchedule(ctx *sweetygo.Context) error {
	scheduleData, err := model.GetSchedule()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(scheduleData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", scheduleData)
}

// GetPartition handles GET /api/:gateway/schedule/partition
func GetPartition(ctx *sweetygo.Context) error {
	partitionData, err := model.GetPartition()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(partitionData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", partitionData)
}

// GetNWStat handles GET /api/:gateway/nwstat
func GetNWStat(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	nwStatData, err := model.GetNWStat(ctx.Param("gateway"), timeRange, now)
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
	timeRange, now := range2stamp(ctx.Param("range"))
	if ctx.Param("advanced") != "" && ctx.Param("advanced") == "1" {
		sensorNWStatAdvData, err := model.GetNWStatAdvByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
		if err != nil {
			return ctx.JSON(500, 0, err.Error(), nil)
		}
		if len(sensorNWStatAdvData) == 0 {
			return ctx.JSON(200, 0, "no result found", nil)
		}
		return ctx.JSON(200, 1, "success", sensorNWStatAdvData)
	}

	sensorNWStatData, err := model.GetNWStatByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(sensorNWStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", sensorNWStatData)
}

// GetLatencyByID handles GET /api/:gateway/nwstat/:id/latency
func GetLatencyByID(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	latency, err := model.GetLatencyByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(latency) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", latency)
}

// GetChInfoByID handles GET /api/:gateway/nwstat/:id/channel
func GetChInfoByID(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	chInfo, err := model.GetChInfoByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(chInfo) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", chInfo)
}

// GetBattery handles GET /api/:gateway/battery
func GetBattery(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	batData, err := model.GetBattery(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(batData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", batData)
}

// GetBatteryByID handles GET /api/:gateway/battery/:sensorID
func GetBatteryByID(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	batData, err := model.GetBatteryByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(batData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", batData)
}

// GetNoiseLevel handles GET /api/:gateway/noise
func GetNoiseLevel(ctx *sweetygo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	nlData, err := model.GetNoiseLevel(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nlData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nlData)
}

// GetTxTotal handles GET /api/:gateway/txtotal
func GetTxTotal(ctx *sweetygo.Context) error {
	timeRange, _ := range2stamp(ctx.Param("range"))
	n, err := model.GetTxTotal(timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if n == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", n)
}

func range2stamp(timeRange string) (int64, int64) {
	now := time.Now().UnixNano() / 1e6
	startTime := int64(0)
	switch timeRange {
	case "15min":
		startTime = now - 15*60*1000
	case "30min":
		startTime = now - 30*60*1000
	case "1hr":
		startTime = now - 60*60*1000
	case "4hr":
		startTime = now - 4*60*60*1000
	case "day":
		startTime = now - 60*60*24*1000
	case "week":
		startTime = now - 60*60*24*7*1000
	case "month":
		startTime = now - 60*60*24*7*30*1000
	default:
		break
	}
	if startTime < lastBootTime {
		startTime = lastBootTime
	}
	return startTime, now
}
