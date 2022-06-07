package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/AmyangXYZ/sgo"
)

// Index page handler.
func Index(ctx *sgo.Context) error {
	return ctx.Render(200, "index")
}

// Static files handler.
func Static(ctx *sgo.Context) error {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
	return nil
}

// GetGateway handles GET /api/gateway
func GetGateway(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	gatewayList, err := modelGetGateway(timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(gatewayList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", gatewayList)
}

// GetTopology handles GET /api/:gateway/topology
func GetTopology(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	nodeList, err := modelGetTopology(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nodeList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nodeList)
}

// GetTopoHistory handles GET /api/:gateway/topology/history
func GetTopoHistory(ctx *sgo.Context) error {
	eventList, err := modelGetTopoHistory(range2stamp("week"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(eventList) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", eventList)
}

// GetSchedule handles GET /api/:gateway/schedule
func GetSchedule(ctx *sgo.Context) error {
	scheduleData, err := modelGetSchedule()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(scheduleData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", scheduleData)
}

// GetPartition handles GET /api/:gateway/schedule/partition
func GetPartition(ctx *sgo.Context) error {
	partitionData, err := modelGetPartition()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(partitionData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", partitionData)
}

// GetPartitionHARP handles GET /api/:gateway/schedule/partition_harp
func GetPartitionHARP(ctx *sgo.Context) error {
	partitionData, err := modelGetPartitionHARP()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(partitionData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", partitionData)
}

// GetSubPartitionHARP handles GET /api/:gateway/schedule/sp_harp
func GetSubPartitionHARP(ctx *sgo.Context) error {
	spData, err := modelGetSubPartitionHARP()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(spData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", spData)
}

// GetNWStat handles GET /api/:gateway/nwstat
func GetNWStat(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	nwStatData, err := modelGetNWStat(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nwStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nwStatData)
}

// GetNWStatByID handles GET /api/:gateway/nwstat/:sensorID
func GetNWStatByID(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	if ctx.Param("advanced") != "" && ctx.Param("advanced") == "1" {
		sensorNWStatAdvData, err := modelGetNWStatAdvByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
		if err != nil {
			return ctx.JSON(500, 0, err.Error(), nil)
		}
		if len(sensorNWStatAdvData) == 0 {
			return ctx.JSON(200, 0, "no result found", nil)
		}
		return ctx.JSON(200, 1, "success", sensorNWStatAdvData)
	}

	sensorNWStatData, err := modelGetNWStatByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(sensorNWStatData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", sensorNWStatData)
}

// GetLatencyByID handles GET /api/:gateway/nwstat/:id/latency
func GetLatencyByID(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	latency, err := modelGetLatencyByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(latency) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", latency)
}

// GetChInfoByID handles GET /api/:gateway/nwstat/:id/channel
func GetChInfoByID(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	chInfo, err := modelGetChInfoByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(chInfo) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", chInfo)
}

// GetBattery handles GET /api/:gateway/battery
func GetBattery(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	batData, err := modelGetBattery(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(batData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", batData)
}

// GetBatteryByID handles GET /api/:gateway/battery/:sensorID
func GetBatteryByID(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	batData, err := modelGetBatteryByID(ctx.Param("gateway"), ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(batData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", batData)
}

// GetNoiseLevel handles GET /api/:gateway/noise
func GetNoiseLevel(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	nlData, err := modelGetNoiseLevel(ctx.Param("gateway"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(nlData) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", nlData)
}

// GetTxTotal handles GET /api/:gateway/txtotal
func GetTxTotal(ctx *sgo.Context) error {
	timeRange, _ := range2stamp(ctx.Param("range"))
	n, err := modelGetTxTotal(timeRange)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if n == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", n)
}

func GetSensorDataByID(ctx *sgo.Context) error {
	timeRange, now := range2stamp(ctx.Param("range"))
	n, err := modelGetSensorDataByID(ctx.Param("sensorID"), timeRange, now)
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(n) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", n)
}

func PutSensorRateByIDnType(ctx *sgo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("sensorID"))
	rate, _ := strconv.ParseFloat(ctx.Param("rate"), 64)
	// err := exec.Command("coap-client", "-m", "put", fmt.Sprintf("[2:%x]/sensors", id), "-e", fmt.Sprintf("%d", int(1/rate))).Run()
	err := exec.Command("curl", fmt.Sprintf("http://host.docker.internal:8080/rate/%d/%d", id, int(1/rate))).Run()
	fmt.Println(fmt.Sprintf("http://host.docker.internal:8080/rate/%d/%d", id, int(1/rate)))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	return ctx.JSON(200, 1, "success", 1)
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
