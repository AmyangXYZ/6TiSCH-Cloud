package model

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func init() {
	dbAddr := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/6tisch", os.Getenv("DBPasswd"))
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
	Timestamp int    `json:"timestamp"`
	Gateway   string `json:"gateway"`
	SensorID  int    `json:"sensor_id"`
	Address   string `json:"address"`
	Parent    int    `json:"parent"`
	Eui64     string `json:"eui64"`
	Position  struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"position"`
	Type  string `json:"type"`
	Power string `json:"power"`
}

func GetGateway(timeRange int64) ([]string, error) {
	var gName string
	gList := make([]string, 0)

	rows, err := db.Query("select distinct GATEWAY_NAME from TOPOLOGY_DATA where LAST_SEEN>=?;", timeRange)
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

	if gatewayName == "any" {
		rows, err = db.Query("select * from TOPOLOGY_DATA where LAST_SEEN>=? group by SENSOR_ID", timeRange)
	} else {
		rows, err = db.Query("select * from TOPOLOGY_DATA where GATEWAY_NAME=? and LAST_SEEN>=? group by SENSOR_ID", gatewayName, timeRange)
	}
	if err != nil {
		return nodeList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&n.LAST_SEEN, &n.Gateway, &n.SensorID,
			&n.Address, &n.Parent, &n.Eui64, &n.Position.Lat, &n.Position.Lng, &n.Type, &n.Power)
		nodeList = append(nodeList, n)
	}
	return nodeList, nil
}

// NWStatData is all sensor's basic network stat data of one gateway
type NWStatData struct {
	SensorID          int     `json:"sensor_id"`
	Gateway           string  `json:"gateway"`
	AVGRTT            float64 `json:"avg_rtt"`
	AvgMACTxTotalDiff float32 `json:"avg_mac_tx_total_diff"`
	AvgMACTxNoACKDiff float32 `json:"avg_mac_tx_noack_diff"`
	AvgAPPPERSentDiff float32 `json:"avg_app_per_sent_diff"`
	AvgAPPPERLostDiff float32 `json:"avg_app_per_lost_diff"`
}

func GetNWStat(gatewayName string, timeRange int64) ([]NWStatData, error) {
	var n NWStatData
	// query NW_DATA_SET_PER_UCONN
	var rows1 *sql.Rows
	// query NW_DATA_SET_LATENCY
	var rows2 *sql.Rows
	nList := make([]NWStatData, 0)

	if gatewayName == "any" {
		rows1, err = db.Query(`select SENSOR_ID, GATEWAY_NAME, AVG(MAC_TX_TOTAL_DIFF),
		AVG(MAC_TX_NOACK_DIFF), AVG(APP_PER_SENT_DIFF),AVG(APP_PER_LOST_DIFF) from NW_DATA_SET_PER_UCONN 
		where TIMESTAMP>=? group by SENSOR_ID`, timeRange)
		if err != nil {
			return nList, err
		}
		rows2, err = db.Query(`select SENSOR_ID, GATEWAY_NAME, AVG(RTT) from
			NW_DATA_SET_LATENCY where TIMESTAMP>=? group by SENSOR_ID`, timeRange)
	} else {
		rows1, err = db.Query(`select SENSOR_ID, GATEWAY_NAME, AVG(MAC_TX_TOTAL_DIFF),
		AVG(MAC_TX_NOACK_DIFF),AVG(APP_PER_SENT_DIFF),AVG(APP_PER_LOST_DIFF) from NW_DATA_SET_PER_UCONN 
		where GATEWAY_NAME=? and TIMESTAMP>=? group by SENSOR_ID`, gatewayName, timeRange)
		if err != nil {
			return nList, err
		}
		rows2, err = db.Query(`select SENSOR_ID, GATEWAY_NAME, AVG(RTT) from
			NW_DATA_SET_LATENCY where GATEWAY_NAME=? and TIMESTAMP>=? group by SENSOR_ID`, gatewayName, timeRange)
	}
	if err != nil {
		return nList, err
	}
	defer rows1.Close()
	// defer rows2.Close()

	for rows1.Next() {
		rows1.Scan(&n.SensorID, &n.Gateway, &n.AvgMACTxTotalDiff, &n.AvgMACTxNoACKDiff,
			&n.AvgAPPPERSentDiff, &n.AvgAPPPERLostDiff)
		nList = append(nList, n)
	}

	// merge
	for rows2.Next() {
		rows2.Scan(&n.SensorID, &n.Gateway, &n.AVGRTT)
		for i, v := range nList {
			if v.SensorID == n.SensorID && v.Gateway == n.Gateway {
				nList[i].AVGRTT = n.AVGRTT
				break
			}
		}
	}

	return nList, nil
}

// SensorNWStatData is each sensor's network statistic: average RSSi value
type SensorNWStatData struct {
	Timestamp int    `json:"timestamp"`
	Gateway   string `json:"gateway"`
	AvgRSSI   int    `json:"avg_rssi"`
}

// SensorNWStatAdvData is each sensor's network statistic detail
type SensorNWStatAdvData struct {
	Timestamp      int    `json:"timestamp"`
	Gateway        string `json:"gateway"`
	MacTxTotalDiff int    `json:"mac_tx_total_diff"`
	MacTxNoAckDiff int    `json:"mac_tx_noack_diff"`
	AppPERSentDiff int    `json:"app_per_sent_diff"`
	AppPERLostDiff int    `json:"app_per_lost_diff"`
}

func GetNWStatByID(gatewayName, sensorID string, timeRange int64) ([]SensorNWStatData, error) {
	var s SensorNWStatData
	var rows *sql.Rows
	sList := make([]SensorNWStatData, 0)

	rows, err = db.Query(`select TIMESTAMP, GATEWAY_NAME, AVG_RSSI from NW_DATA_SET_PER_UCONN 
			where GATEWAY_NAME=? and SENSOR_ID=? and TIMESTAMP>=?`, gatewayName, sensorID, timeRange)
	if err != nil {
		return sList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&s.Timestamp, &s.Gateway, &s.AvgRSSI)
		sList = append(sList, s)
	}

	return sList, nil
}

func GetNWStatAdvByID(gatewayName, sensorID string, timeRange int64) ([]SensorNWStatAdvData, error) {
	var s SensorNWStatAdvData
	var rows *sql.Rows
	sList := make([]SensorNWStatAdvData, 0)

	rows, err = db.Query(`select TIMESTAMP, GATEWAY_NAME, MAC_TX_TOTAL_DIFF,
			MAC_TX_NOACK_DIFF,APP_PER_SENT_DIFF,APP_PER_LOST_DIFF from NW_DATA_SET_PER_UCONN 
			where GATEWAY_NAME=? and SENSOR_ID=? and TIMESTAMP>=?`, gatewayName, sensorID, timeRange)
	if err != nil {
		return sList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&s.Timestamp, &s.Gateway, &s.MacTxTotalDiff, &s.MacTxNoAckDiff,
			&s.AppPERSentDiff, &s.AppPERLostDiff)
		sList = append(sList, s)
	}

	return sList, nil
}

type SensorBatteryData struct {
	Gateway         string  `json:"gateway"`
	SensorID        int     `json:"sensor_id"`
	AvgCC2650Active float64 `json:"avg_cc2650_active"`
	AvgCC2650Sleep  float64 `json:"avg_cc2650_sleep"`
	AvgRFRx         float64 `json:"avg_rf_rx"`
	AvgRFTx         float64 `json:"avg_rf_tx"`
	BatRemain       string  `json:"bat_remain"`
}

func GetBattery(gatewayName string, timeRange int64) ([]SensorBatteryData, error) {
	var b SensorBatteryData
	var rows *sql.Rows
	bList := make([]SensorBatteryData, 0)
	var bat float64
	if gatewayName == "any" {
		rows, err = db.Query(`select SENSOR_ID,GATEWAY_NAME,AVG(CC2650_ACTIVE),AVG(CC2650_SLEEP),AVG(RF_RX),AVG(RF_TX),BAT 
			from SENSOR_DATA where TIMESTAMP>=? group by SENSOR_ID`, timeRange)
	} else {
		rows, err = db.Query(`select SENSOR_ID,GATEWAY_NAME,AVG(CC2650_ACTIVE),AVG(CC2650_SLEEP),AVG(RF_RX),AVG(RF_TX),BAT 
			from SENSOR_DATA where GATEWAY_NAME=? and TIMESTAMP>=? group by SENSOR_ID`, gatewayName, timeRange)
	}
	if err != nil {
		return bList, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&b.SensorID, &b.Gateway, &b.AvgCC2650Active, &b.AvgCC2650Sleep, &b.AvgRFRx, &b.AvgRFTx, &bat)
		// todo
		// b.BatRemain = string(bat)
		bList = append(bList, b)
	}
	return bList, nil
}
