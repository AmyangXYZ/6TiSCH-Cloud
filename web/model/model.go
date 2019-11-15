package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	dbAddr = "root:1234@tcp(127.0.0.1:3306)/6tisch"
	err    error
)

func init() {
	db, _ = sql.Open("mysql", dbAddr)
	for {
		if err := db.Ping(); err != nil {
			fmt.Println(err, ", retry in 10s...")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	// https://github.com/go-sql-driver/mysql/issues/674
	db.SetMaxIdleConns(0)
}

// Node info for topology.
type Node struct {
	Datetime    string `json:"datetime"`
	Timestamp   int    `json:"timestamp"`
	GatewayName string `json:"gateway_name"`
	SensorID    int    `json:"sensor_id"`
	Address     string `json:"address"`
	Parent      int    `json:"parent"`
	Eui64       string `json:"eui64"`
	Position    struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"position"`
	Type  string `json:"type"`
	Power string `json:"power"`
}

func GetGateway(timeRange int64) ([]string, error) {
	var gName string
	gList := make([]string, 0)

	rows, err := db.Query("select distinct GATEWAY_NAME from TOPOLOGY_DATA where TIMESTAMP>?;", timeRange)
	if err != nil {
		return gList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&gName)
		gList = append(gList, gName)
	}
	// for multi-gateway test
	gList = append(gList, "UCONN_GWX")
	return gList, nil
}

func GetTopology(gatewayName string, timeRange int64) ([]Node, error) {
	var n Node
	var rows *sql.Rows
	nodeList := make([]Node, 0)

	if gatewayName == "ANY" {
		rows, err = db.Query("select * from TOPOLOGY_DATA where TIMESTAMP>?", timeRange)
	} else {
		rows, err = db.Query("select * from TOPOLOGY_DATA where GATEWAY_NAME=? and TIMESTAMP>?", gatewayName, timeRange)
	}
	if err != nil {
		return nodeList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&n.Datetime, &n.Timestamp, &n.GatewayName, &n.SensorID,
			&n.Address, &n.Parent, &n.Eui64, &n.Position.Lat, &n.Position.Lng, &n.Type, &n.Power)
		nodeList = append(nodeList, n)
	}
	return nodeList, nil
}

// NWStatData is all sensor's basic network stat data of one gateway
type NWStatData struct {
	SensorID          int     `json:"sensor_id"`
	AVGRTT            float64 `json:"avg_rtt"`
	AvgMACTxTotalDiff float32 `json:"avg_mac_tx_total_diff"`
	AvgMACTxNoACKDiff float32 `json:"avg_mac_tx_noack_diff"`
	AvgAPPPERSentDiff float32 `json:"avg_app_per_sent_diff"`
	AvgAPPPERLostDiff float32 `json:"avg_app_per_lost_diff"`
}

func GetNWStat(gatewayName string, timeRange int64) ([]NWStatData, error) {
	var n NWStatData
	var rows *sql.Rows
	nList := make([]NWStatData, 0)

	if gatewayName == "ANY" {
		rows, err = db.Query(`select NW_DATA_SET_LATENCY.SENSOR_ID, AVG(RTT), 
		AVG(MAC_TX_TOTAL_DIFF),AVG(MAC_TX_NOACK_DIFF),AVG(APP_PER_SENT_DIFF),AVG(APP_PER_LOST_DIFF) 
		from NW_DATA_SET_LATENCY inner join NW_DATA_SET_PER_UCONN on NW_DATA_SET_LATENCY.TIMESTAMP>? and
		NW_DATA_SET_LATENCY.SENSOR_ID=NW_DATA_SET_PER_UCONN.SENSOR_ID group by SENSOR_ID`, timeRange)
	} else {
		rows, err = db.Query(`select NW_DATA_SET_LATENCY.SENSOR_ID, AVG(RTT), 
		AVG(MAC_TX_TOTAL_DIFF),AVG(MAC_TX_NOACK_DIFF),AVG(APP_PER_SENT_DIFF),AVG(APP_PER_LOST_DIFF) 
		from NW_DATA_SET_LATENCY inner join NW_DATA_SET_PER_UCONN on NW_DATA_SET_LATENCY.GATEWAY_NAME=? and
		NW_DATA_SET_LATENCY.TIMESTAMP>? and NW_DATA_SET_LATENCY.SENSOR_ID=NW_DATA_SET_PER_UCONN.SENSOR_ID group by SENSOR_ID`, gatewayName, timeRange)
	}
	if err != nil {
		return nList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&n.SensorID, &n.AVGRTT, &n.AvgMACTxTotalDiff,
			&n.AvgMACTxNoACKDiff, &n.AvgAPPPERSentDiff,
			&n.AvgAPPPERLostDiff)
		nList = append(nList, n)
	}

	return nList, nil
}

// SensorNWStatData is each sensor's network statistic detail
type SensorNWStatData struct {
	Timestamp int `json:"timestamp"`
	Ch        map[string]struct {
		RSSI    int `json:"rssi"`
		RxRSSI  int `json:"rx_rssi"`
		TxNoAck int `json:"tx_noack"`
		TxTotal int `json:"tx_total"`
	} `json:"ch"`
	MacTxTotalDiff int `json:"mac_tx_total_diff"`
	MacTxNoAckDiff int `json:"mac_tx_noack_diff"`
	AppPERSentDiff int `json:"app_per_sent_diff"`
	AppPERLostDiff int `json:"app_per_lost_diff"`
}

func GetNWStatByID(gatewayName, sensorID string, timeRange int64) ([]SensorNWStatData, error) {
	var s SensorNWStatData
	var rows *sql.Rows
	var chInfo string
	sList := make([]SensorNWStatData, 0)

	rows, err = db.Query(`select TIMESTAMP,CHANNEL_INFO,MAC_TX_TOTAL_DIFF,
		MAC_TX_NOACK_DIFF,APP_PER_SENT_DIFF,APP_PER_LOST_DIFF from NW_DATA_SET_PER_UCONN 
		where GATEWAY_NAME=? and SENSOR_ID=? and TIMESTAMP>`, gatewayName, sensorID, timeRange)
	if err != nil {
		return sList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&s.Timestamp, &chInfo, &s.MacTxTotalDiff, &s.MacTxNoAckDiff,
			&s.AppPERSentDiff, &s.AppPERLostDiff)
		err = json.Unmarshal([]byte(chInfo), &s.Ch)
		if err != nil {
			return sList, err
		}
		sList = append(sList, s)
	}
	return sList, nil
}
