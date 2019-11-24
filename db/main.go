package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Info        *log.Logger
	Error       *log.Logger
	infoHandle  = os.Stdout
	errorHandle = os.Stdout

	db *sql.DB
)

func main() {
	http.HandleFunc("/biu", handle)
	http.ListenAndServe(":54321", nil)
}

// parse and store json data sent from gateway
func handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(err)
		return
	}

	var msg json.RawMessage
	d := data{Msg: &msg}
	err = json.Unmarshal(body, &d)
	if err != nil {
		Error.Println(err)
		return
	}
	Info.Println(d.Type)

	gwn := d.Gateway.Msg.Name

	switch d.Type {
	case "heart":
		var h heart
		err = json.Unmarshal(msg, &h)
		if err != nil {
			Error.Println(err)
			return
		}
		handleHeartBeatData(h)
	case "topology_data":
		var t topology
		err = json.Unmarshal(msg, &t)
		if err != nil {
			Error.Println(err)
			return
		}
		handleTopologyData(t, gwn)
	case "nodes_data":
		var n []node
		err = json.Unmarshal(msg, &n)
		if err != nil {
			Error.Println(err)
			return
		}
		handleNodesData(n, gwn)
	case "sensor_type_0":
		var s sensor
		err = json.Unmarshal(msg, &s)
		if err != nil {
			Error.Println(err)
			return
		}
		handleSensorData(s, gwn)
	case "network_data_0":
		var n0 network0
		err = json.Unmarshal(msg, &n0)
		if err != nil {
			Error.Println(err)
			return
		}
		handleNetworkData0(n0, gwn)
	case "network_data_1":
		var n1 network1
		err = json.Unmarshal(msg, &n1)
		if err != nil {
			Error.Println(err)
			return
		}
		handleNetworkData1(n1, gwn)
	case "network_data_2":
		var n2 network2
		err = json.Unmarshal(msg, &n2)
		if err != nil {
			Error.Println(err)
			return
		}
		handleNetworkData2(n2, gwn)
	default:
		Error.Println("Unknown data type:", string(body))
	}

	fmt.Fprintf(w, "Got it!\n")
}

func handleTopologyData(topo topology, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt, err := db.Prepare(`INSERT INTO TOPOLOGY_DATA(FIRST_APPEAR, LAST_SEEN, GATEWAY_NAME, SENSOR_ID, ADDRESS,
		PARENT, EUI64, GPS_Lat, GPS_Lon, TYPE, POWER) VALUES(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}

	_, err = stmt.Exec(timestamp, timestamp, gwn, topo.Data.ID, topo.Data.Address,
		topo.Data.Parent, topo.Data.Eui64, topo.Data.GPS[0], topo.Data.GPS[1], topo.Data.Type, topo.Data.Power)
	if err != nil {
		Error.Println(err)
	}
}

func handleHeartBeatData(h heart) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt, err := db.Prepare(`UPDATE TOPOLOGY_DATA SET LAST_SEEN=? where GATEWAY_NAME=? and SENSOR_ID=1`)
	if err != nil {
		Error.Println(err)
	}
	_, err = stmt.Exec(timestamp, h.Msg.Name)
	if err != nil {
		Error.Println(err)
	}
}

func handleNodesData(n []node, gwn string) {
	// fmt.Println(n)
}
func handleSensorData(s sensor, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt1, err := db.Prepare(`INSERT INTO SENSOR_DATA (TIMESTAMP, GATEWAY_NAME, SENSOR_ID, TEMP, 
		RHUM, LUX, PRESS, ACCELX, ACCELY, ACCELZ, LED, EH, EH1, CC2650_ACTIVE, CC2650_SLEEP,RF_TX, RF_RX, 
		MSP432_ACTIVE, MSP432_SLEEP, GPSEN_ACTIVE, GPSEN_SLEEP, OTHERS, SEQUENCE, ASN_STAMP1, ASN_STAMP2, CHANNEL, BAT, LATENCY) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt1, err)
	}

	_, err = stmt1.Exec(timestamp, gwn, s.ID, s.Data.Temp, s.Data.Rhum, s.Data.Lux, s.Data.Press,
		s.Data.Accelx, s.Data.Accely, s.Data.Accelz, s.Data.LED, s.Data.Eh, s.Data.Eh1, s.Data.CC2650Active, s.Data.CC2650Sleep,
		s.Data.RFTx, s.Data.RFRx, s.Data.MSP432Active, s.Data.MSP432Sleep, s.Data.GPSEnActive, s.Data.GPSEnSleep, s.Data.Others,
		s.Data.Sequence, s.Data.ASNStamp1, s.Data.ASNStamp2, s.Data.Channel, s.Data.Bat, s.Data.Latency)
	if err != nil {
		Error.Println(stmt1, err)
	}

	stmt2, err := db.Prepare(`UPDATE TOPOLOGY_DATA SET LAST_SEEN=? where GATEWAY_NAME=? and SENSOR_ID=?`)
	if err != nil {
		Error.Println(stmt2, err)
	}
	_, err = stmt2.Exec(timestamp, gwn, s.ID)
	if err != nil {
		Error.Println(stmt2, err)
	}
}
func handleNetworkData0(n0 network0, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt1, err := db.Prepare(`INSERT INTO NW_DATA_SET_PER_UCONN(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, 
		AVG_RSSI, APP_PER_SENT_LAST_SEQ, APP_PER_SENT, APP_PER_SENT_LOST, TX_FAIL, TX_NOACK, 
		TX_TOTAL, RX_TOTAL, TX_LENGTH_TOTAL, MAC_TX_NOACK_DIFF, MAC_TX_TOTAL_DIFF, MAC_RX_TOTAL_DIFF, 
		MAC_TX_LENGTH_TOTAL_DIFF, APP_PER_LOST_DIFF, APP_PER_SENT_DIFF) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt1, err)
	}
	stmt2, err := db.Prepare(`INSERT INTO NW_DATA_SET_PER_CHINFO(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, 
		CHANNELS, RSSI, RX_RSSI, TX_NOACK, TX_TOTAL) 
		VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt2, err)
	}
	// compute average rssi and store avaliable channel info into NW_DATA_SET_PER_CHINFO
	// save 95% space than in text field
	avgRSSi := 0
	cnt := 0
	tmp := 0

	chList := ""
	rssiList := ""
	rxRSSiList := ""
	txNoAckList := ""
	txTotalList := ""

	for k, v := range n0.Data.Ch {
		if v.RSSI != 0 {
			cnt++
			tmp += v.RSSI

			chList += k + ","
			rssiList += strconv.Itoa(v.RSSI) + ","
			rxRSSiList += strconv.Itoa(v.RxRSSI) + ","
			txNoAckList += strconv.Itoa(v.TxNoAck) + ","
			txTotalList += strconv.Itoa(v.TxTotal) + ","
		}
	}
	avgRSSi = tmp / cnt

	_, err = stmt1.Exec(timestamp, gwn, n0.ID, avgRSSi, n0.Data.AppPER.LastSeq,
		n0.Data.AppPER.Sent, n0.Data.AppPER.Lost, n0.Data.TxFail, n0.Data.TxNoAck, n0.Data.TxTotal,
		n0.Data.RxTotal, n0.Data.TxLengthTotal, n0.Data.MacTxNoAckDiff, n0.Data.MacTxTotalDiff,
		n0.Data.MacRxTotalDiff, n0.Data.MacTxLengthTotalDiff, n0.Data.AppLostDiff, n0.Data.AppSentDiff)
	if err != nil {
		Error.Println(stmt1, err)
	}
	_, err = stmt2.Exec(timestamp, gwn, n0.ID, chList, rssiList, rxRSSiList, txNoAckList, txTotalList)
	if err != nil {
		Error.Println(stmt2, err)
	}
}

func handleNetworkData1(n1 network1, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt, err := db.Prepare(`INSERT INTO NW_DATA_SET_INFO(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, 
		CUR_PARENT, NUM_PARENT_CHANGE, NUM_SYNC_LOST, AVG_DRIFT, MAX_DRIFT, NUM_MAC_OUT_OF_BUFFER, NUM_UIP_RX_LOST, 
		NUM_LOWPAN_TX_LOST, NUM_LOWPAN_RX_LOST, NUM_COAP_RX_LOST, NUM_COAP_OBS_DIS) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}

	_, err = stmt.Exec(timestamp, gwn, n1.ID, n1.Data.CurParent, n1.Data.NumParentChange,
		n1.Data.NumSyncLost, n1.Data.AvgDrift, n1.Data.MaxDrift, n1.Data.NumMacOutOfBuffer, n1.Data.NumUipRxLost,
		n1.Data.NumLowpanTxLost, n1.Data.NumLowpanRxLost, n1.Data.NumCoapRxLost, n1.Data.NumCoapObsDis)
	if err != nil {
		Error.Println(stmt, err)
	}
}

func handleNetworkData2(n2 network2, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt, err := db.Prepare(`INSERT INTO NW_DATA_SET_LATENCY(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, RTT) 
		VALUES(?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}

	_, err = stmt.Exec(timestamp, gwn, n2.ID, n2.RTT)
	if err != nil {
		Error.Println(stmt, err)
	}
}

// init logger and db
func init() {
	Info = log.New(infoHandle, "[*] INFO: ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "[!] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	dbAddr := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/6tisch", os.Getenv("DBPasswd"))
	db, _ = sql.Open("mysql", dbAddr)
	for {
		if err := db.Ping(); err != nil {
			Error.Println(err, ", retry in 10s...")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	Info.Println("connected to db")
	// https://github.com/go-sql-driver/mysql/issues/674
	db.SetMaxIdleConns(0)

	// TOPOLOGY_DATA
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS TOPOLOGY_DATA (
			FIRST_APPEAR BIGINT,
			LAST_SEEN BIGINT,
			GATEWAY_NAME VARCHAR(16) NOT NULL,
			SENSOR_ID SMALLINT UNSIGNED NOT NULL,
			ADDRESS VARCHAR(64) NOT NULL,
			PARENT SMALLINT,
			EUI64 VARCHAR(64),
			GPS_Lat DOUBLE NOT NULL,
			GPS_Lon DOUBLE NOT NULL,
			TYPE VARCHAR(64) NOT NULL,
			POWER VARCHAR(64) NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table TOPOLOGY_DATA ready")

	// SENSOR_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS SENSOR_DATA (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		TEMP SMALLINT UNSIGNED NOT NULL,
		RHUM SMALLINT UNSIGNED NOT NULL,
		LUX SMALLINT UNSIGNED NOT NULL,
		PRESS SMALLINT UNSIGNED NOT NULL,
		ACCELX FLOAT NOT NULL,
		ACCELY FLOAT NOT NULL,
		ACCELZ FLOAT NOT NULL,
		LED SMALLINT UNSIGNED NOT NULL,
		EH SMALLINT UNSIGNED NOT NULL,
		EH1 SMALLINT UNSIGNED NOT NULL,
		CC2650_ACTIVE TINYINT UNSIGNED NOT NULL,
		CC2650_SLEEP TINYINT UNSIGNED NOT NULL,
		RF_TX FLOAT NOT NULL,
		RF_RX FLOAT NOT NULL,
		MSP432_ACTIVE TINYINT UNSIGNED NOT NULL,
		MSP432_SLEEP TINYINT UNSIGNED NOT NULL,
		GPSEN_ACTIVE FLOAT NOT NULL,
		GPSEN_SLEEP FLOAT NOT NULL,
		OTHERS SMALLINT NOT NULL,
		SEQUENCE SMALLINT UNSIGNED NOT NULL,
		ASN_STAMP1 SMALLINT NOT NULL,
		ASN_STAMP2 SMALLINT NOT NULL,
		CHANNEL TINYINT UNSIGNED NOT NULL,
		BAT FLOAT NOT NULL,
		LATENCY FLOAT NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table SENSOR_DATA ready")

	// NW_DATA_SET_PER_UCONN or network_data_0
	db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_PER_UCONN (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		CHANNEL_INFO TEXT,
		AVG_RSSI SMALLINT NOT NULL,
		APP_PER_SENT_LAST_SEQ SMALLINT NOT NULL,
		APP_PER_SENT SMALLINT NOT NULL,
		APP_PER_SENT_LOST SMALLINT NOT NULL,
		TX_FAIL INT NOT NULL,
		TX_NOACK INT NOT NULL,
		TX_TOTAL INT NOT NULL,
		RX_TOTAL INT NOT NULL,
		TX_LENGTH_TOTAL INT NOT NULL,
		MAC_TX_NOACK_DIFF SMALLINT NOT NULL,
		MAC_TX_TOTAL_DIFF SMALLINT NOT NULL,
		MAC_RX_TOTAL_DIFF SMALLINT NOT NULL,
		MAC_TX_LENGTH_TOTAL_DIFF SMALLINT NOT NULL,
		APP_PER_LOST_DIFF SMALLINT NOT NULL,
		APP_PER_SENT_DIFF SMALLINT NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table NW_DATA_SET_PER_UCONN ready")

	// CHANNEL_INFO in NW_DATA_SET_PER_UCONN or network_data_0
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_PER_CHINFO (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		CHANNELS VARCHAR(64) NOT NULL,
		RSSI VARCHAR(128) NOT NULL,
		RX_RSSI  VARCHAR(128) NOT NULL,
		TX_NOACK VARCHAR(128) NOT NULL,
		TX_TOTAL VARCHAR(64) NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table NW_DATA_SET_PER_CHINFO ready")

	// NW_DATA_SET_INFO or network_data_1
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_INFO (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		CUR_PARENT SMALLINT UNSIGNED NOT NULL,
		NUM_PARENT_CHANGE SMALLINT UNSIGNED NOT NULL,
		NUM_SYNC_LOST SMALLINT UNSIGNED NOT NULL,
		AVG_DRIFT SMALLINT UNSIGNED NOT NULL,
		MAX_DRIFT SMALLINT UNSIGNED NOT NULL,
		NUM_MAC_OUT_OF_BUFFER SMALLINT UNSIGNED NOT NULL,
		NUM_UIP_RX_LOST SMALLINT UNSIGNED NOT NULL,
		NUM_LOWPAN_TX_LOST SMALLINT UNSIGNED NOT NULL,
		NUM_LOWPAN_RX_LOST SMALLINT UNSIGNED NOT NULL,
		NUM_COAP_RX_LOST SMALLINT UNSIGNED NOT NULL,
		NUM_COAP_OBS_DIS SMALLINT UNSIGNED NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table NW_DATA_SET_INFO ready")

	// NW_DATA_SET_LATENCY or network_data_2
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_LATENCY (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		RTT FLOAT NOT NULL);`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table NW_DATA_SET_LATENCY ready")
	Info.Println("db ready")
}

// JSON data format
type (
	data struct {
		Type    string      `json:"type"`
		Gateway gateway     `json:"gateway_0"`
		Msg     interface{} `json:"msg"`
	}

	gateway struct {
		Type string `json:"type"`
		Msg  struct {
			Name    string     `json:"name"`
			Address string     `json:"address"`
			GPS     [2]float64 `json:"gps"`
			Color   string     `json:"color"`
		} `json:"msg"`
	}

	heart struct {
		Type string `json:"type"`
		Msg  struct {
			Name    string     `json:"name"`
			Address string     `json:"address"`
			GPS     [2]float64 `json:"gps"`
			Color   string     `json:"color"`
		} `json:"msg"`
	}

	node struct {
		ID          int        `json:"_id"`
		Address     string     `json:"address"`
		Eui64       string     `json:"eui64"`
		Candidate   []int      `json:"candidate"`
		LifeTime    int        `json:"lifetime"`
		Capacity    int        `json:"capacity"`
		BeaconState string     `json:"beacon_state"`
		Sensors     dataSensor `json:"sensors"`
		AppPer      int        `json:"app_per"`
		Meta        struct {
			GPS   [2]float64 `json:"gps"`
			Type  string     `json:"type"`
			Power string     `json:"power"`
		} `json:"meta"`
	}

	topology struct {
		Data dataTopology `json:"data"`
	}
	dataTopology struct {
		ID      int        `json:"_id"`
		Address string     `json:"address"`
		Parent  int        `json:"parent"`
		Eui64   string     `json:"eui64"`
		GPS     [2]float64 `json:"gps"`
		Type    string     `json:"type"`
		Power   string     `json:"power"`
	}

	sensor struct {
		ID   int        `json:"_id"`
		Data dataSensor `json:"data"`
	}
	dataSensor struct {
		Temp         float32 `json:"temp"`
		Rhum         int     `json:"rhum"`
		Lux          int     `json:"lux"`
		Press        int     `json:"press"`
		Accelx       float32 `json:"accelx"`
		Accely       float32 `json:"accely"`
		Accelz       float32 `json:"type"`
		LED          int     `json:"led"`
		Eh           int     `json:"eh"`
		Eh1          int     `json:"eh1"`
		CC2650Active int     `json:"cc2650_active"`
		CC2650Sleep  int     `json:"cc2650_sleep"`
		RFTx         float32 `json:"rf_tx"`
		RFRx         float32 `json:"rf_rx"`
		MSP432Active int     `json:"msp432_active"`
		MSP432Sleep  int     `json:"msp432_sleep"`
		GPSEnActive  float32 `json:"gpsen_active"`
		GPSEnSleep   float32 `json:"gpsen_sleep"`
		Others       int     `json:"others"`
		Sequence     int     `json:"sequence"`
		ASNStamp1    int     `json:"asn_stamp1"`
		ASNStamp2    int     `json:"asn_stamp2"`
		Channel      int     `json:"channel"`
		Bat          float32 `json:"bat"`
		Latency      float32 `json:"latency"`
	}

	network0 struct {
		ID   int          `json:"_id"`
		Data networkData0 `json:"data"`
	}
	network1 struct {
		ID   int          `json:"_id"`
		Data networkData1 `json:"data"`
	}
	network2 struct {
		ID   int          `json:"_id"`
		RTT  float32      `json:"rtt"`
		Data networkData2 `json:"data"`
	}

	networkData0 struct {
		Ch map[string]struct {
			RSSI    int `json:"rssi"`
			RxRSSI  int `json:"rxRssi"`
			TxNoAck int `json:"txNoAck"`
			TxTotal int `json:"txTotal"`
		} `json:"ch"`
		AppPER struct {
			LastSeq int `json:"last_seq"`
			Sent    int `json:"sent"`
			Lost    int `json:"lost"`
		} `json:"appPer"`
		TxFail               int `json:"txFail"`
		TxNoAck              int `json:"txNoAck"`
		TxTotal              int `json:"txTotal"`
		RxTotal              int `json:"rxTotal"`
		TxLengthTotal        int `json:"txLengthTotal"`
		MacTxNoAckDiff       int `json:"macTxNoAckDiff"`
		MacTxTotalDiff       int `json:"macTxTotalDiff"`
		MacRxTotalDiff       int `json:"macRxTotalDiff"`
		MacTxLengthTotalDiff int `json:"macTxLengthTotalDiff"`
		AppLostDiff          int `json:"appLostDiff"`
		AppSentDiff          int `json:"appSentDiff"`
	}

	networkData1 struct {
		CurParent         int `json:"curParent"`
		NumParentChange   int `json:"numParentChange"`
		NumSyncLost       int `json:"numSyncLost"`
		AvgDrift          int `json:"avgDrift"`
		MaxDrift          int `json:"maxDrift"`
		NumMacOutOfBuffer int `json:"numMacOutOfBuffer"`
		NumUipRxLost      int `json:"numUipRxLost"`
		NumLowpanTxLost   int `json:"numLowpanTxLost"`
		NumLowpanRxLost   int `json:"numLowpanRxLost"`
		NumCoapRxLost     int `json:"numCoapRxLost"`
		NumCoapObsDis     int `json:"numCoapObsDis"`
	}

	networkData2 struct{}
)
