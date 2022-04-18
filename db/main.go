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

	go func() {

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
		case "partition_data":
			var pt partition
			err = json.Unmarshal(msg, &pt)
			if err != nil {
				Error.Println(err)
				return
			}
			handlePartitionData(pt, gwn)
		case "harp_partition_data":
			var pt partition
			err = json.Unmarshal(msg, &pt)
			if err != nil {
				Error.Println(err)
				return
			}
			handlePartitionHARPData(pt, gwn)
		case "harp_sp_data":
			var sp subpartition
			err = json.Unmarshal(msg, &sp)
			if err != nil {
				Error.Println(err)
				return
			}
			handleHARPSPData(sp, gwn)
		case "schedule_data":
			var sch schedule
			err = json.Unmarshal(msg, &sch)
			if err != nil {
				Error.Println(err)
				return
			}
			handleScheduleData(sch, gwn)
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

	}()
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
	defer stmt.Close()
	_, err = stmt.Exec(timestamp, timestamp, gwn, topo.Data.ID, topo.Data.Address,
		topo.Data.Parent, topo.Data.Eui64, topo.Data.GPS[0], topo.Data.GPS[1], topo.Data.Type, topo.Data.Power)
	if err != nil {
		Error.Println(err)
	}
}

func handleHeartBeatData(h heart) {
}

func handlePartitionData(pt partition, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	// clear table
	_, err := db.Exec("DELETE FROM PARTITION_DATA where 1")
	if err != nil {
		Error.Println(err)
	}

	stmt, err := db.Prepare(`INSERT INTO PARTITION_DATA(TIMESTAMP, GATEWAY_NAME, ROWW, TYPE, LAYER, CHANNEL_START,CHANNEL_END, START, END) VALUES(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}
	defer stmt.Close()
	for i := range pt.Data {
		_, err = stmt.Exec(timestamp, gwn, pt.Data[i].Row, pt.Data[i].Type, pt.Data[i].Layer, pt.Data[i].Channels[0], pt.Data[i].Channels[1], pt.Data[i].Range[0], pt.Data[i].Range[1])
		if err != nil {
			Error.Println(err)
		}
	}
}

func handlePartitionHARPData(pt partition, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	// clear table
	_, err := db.Exec("DELETE FROM HARP_PARTITION_DATA where 1")
	if err != nil {
		Error.Println(err)
	}
	// clear table
	_, err = db.Exec("DELETE FROM HARP_SP_DATA where 1")
	if err != nil {
		Error.Println(err)
	}
	// clear table
	_, err = db.Exec("DELETE FROM SCHEDULE_DATA where 1")
	if err != nil {
		Error.Println(err)
	}

	stmt, err := db.Prepare(`INSERT INTO HARP_PARTITION_DATA(TIMESTAMP, GATEWAY_NAME, TYPE, START, END) VALUES(?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}
	defer stmt.Close()
	for i := range pt.Data {
		_, err = stmt.Exec(timestamp, gwn, pt.Data[i].Type, pt.Data[i].Range[0], pt.Data[i].Range[1])
		if err != nil {
			Error.Println(err)
		}
	}
}

func handleHARPSPData(sp subpartition, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt0, err := db.Prepare(`DELETE FROM HARP_SP_DATA where SENSOR_ID=?`)
	if err != nil {
		Error.Println(stmt0, err)
	}
	defer stmt0.Close()

	_, err = stmt0.Exec(sp.ID)
	if err != nil {
		Error.Println(err)
	}

	stmt, err := db.Prepare(`INSERT INTO HARP_SP_DATA (TIMESTAMP, GATEWAY_NAME, SENSOR_ID, LAYER, TS_START, TS_END, CH_START, CH_END) 
                VALUES (?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}
	defer stmt.Close()
	for i := range sp.Data {
		_, err = stmt.Exec(timestamp, gwn, sp.ID, sp.Data[i].Layer, sp.Data[i].TsStart, sp.Data[i].TsEnd, sp.Data[i].ChStart, sp.Data[i].ChEnd)
		if err != nil {
			Error.Println(err)
		}
	}
}

func handleScheduleData(sch schedule, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	// clear table
	_, err := db.Exec("DELETE FROM SCHEDULE_DATA where 1")
	if err != nil {
		Error.Println(err)
	}

	stmt, err := db.Prepare(`INSERT INTO SCHEDULE_DATA(TIMESTAMP, GATEWAY_NAME, SLOT_OFFSET, CHANNEL_OFFSET, SUBSLOT_OFFSET, SUBSLOT_PERIOD, TYPE, LAYER, SENDER, RECEIVER, IS_OPTIMAL) VALUES(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}
	defer stmt.Close()
	for i := range sch.Data {
		_, err = stmt.Exec(timestamp, gwn, sch.Data[i].Slot[0], sch.Data[i].Slot[1],
			sch.Data[i].Subslot[0], sch.Data[i].Subslot[1], sch.Data[i].Cell.Type, sch.Data[i].Cell.Layer, sch.Data[i].Cell.Sender, sch.Data[i].Cell.Receiver, sch.Data[i].IsOptimal)
		if err != nil {
			Error.Println(err)
		}
	}
}

func handleSensorData(s sensor, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt1, err := db.Prepare(`INSERT INTO SENSOR_DATA (TIMESTAMP, GATEWAY_NAME, SENSOR_ID, TEMP, HUMIDITY,ULTRASONIC,LVDT1,LVDT2,
                RHUM, LUX, PRESS, ACCELX, ACCELY, ACCELZ, LED, EH, EH1, CC2650_ACTIVE, CC2650_SLEEP,RF_TX, RF_RX, 
                MSP432_ACTIVE, MSP432_SLEEP, GPSEN_ACTIVE, GPSEN_SLEEP, OTHERS, SEQUENCE, ASN_STAMP1, CHANNEL, BAT, A2A_LATENCY, LAST_UPLINK_LATENCY) 
                VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt1, err)
	}
	defer stmt1.Close()
	_, err = stmt1.Exec(timestamp, gwn, s.ID, s.Data.Temp, s.Data.Humidity, s.Data.Ultrasonic, s.Data.LVDT1, s.Data.LVDT2, s.Data.Rhum, s.Data.Lux, s.Data.Press,
		s.Data.Accelx, s.Data.Accely, s.Data.Accelz, s.Data.LED, s.Data.Eh, s.Data.Eh1, s.Data.CC2650Active, s.Data.CC2650Sleep,
		s.Data.RFTx, s.Data.RFRx, s.Data.MSP432Active, s.Data.MSP432Sleep, s.Data.GPSEnActive, s.Data.GPSEnSleep, s.Data.Others,
		s.Data.Sequence, s.Data.ASNStamp1, s.Data.Channel, s.Data.Bat, s.Data.A2ALatency, s.Data.LastUplinkLatency)
	if err != nil {
		Error.Println(stmt1, err)
	}

	stmt2, err := db.Prepare(`UPDATE TOPOLOGY_DATA SET LAST_SEEN=? where (GATEWAY_NAME=? and SENSOR_ID=?) or SENSOR_ID=1`)
	if err != nil {
		Error.Println(stmt2, err)
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(timestamp, gwn, s.ID)
	if err != nil {
		Error.Println(stmt2, err)
	}
}
func handleNetworkData0(n0 network0, gwn string) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6

	stmt1, err := db.Prepare(`INSERT INTO NW_DATA_SET_PER_UCONN(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, 
                AVG_RSSI, AVG_RXRSSI, APP_PER_SENT_LAST_SEQ, APP_PER_SENT, APP_PER_SENT_LOST, TX_FAIL, TX_NOACK, 
                TX_TOTAL, RX_TOTAL, TX_LENGTH_TOTAL, MAC_TX_NOACK_DIFF, MAC_TX_TOTAL_DIFF, MAC_RX_TOTAL_DIFF, 
                MAC_TX_LENGTH_TOTAL_DIFF, APP_PER_LOST_DIFF, APP_PER_SENT_DIFF) 
                VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt1, err)
	}
	defer stmt1.Close()
	stmt2, err := db.Prepare(`INSERT INTO NW_DATA_SET_PER_CHINFO(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, 
                CHANNELS, RSSI, RX_RSSI, TX_NOACK, TX_TOTAL) 
                VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		Error.Println(stmt2, err)
	}
	defer stmt2.Close()
	// compute average rssi,rxrssi and store avaliable channel info into NW_DATA_SET_PER_CHINFO
	// save 95% space than in text field
	cnt := 0
	tmp := 0
	tmp2 := 0
	chList := ""
	rssiList := ""
	rxRSSiList := ""
	txNoAckList := ""
	txTotalList := ""

	for k, v := range n0.Data.Ch {
		if v.RSSI != 0 {
			cnt++
			tmp += v.RSSI
			tmp2 += v.RxRSSI
			chList += k + ","
			rssiList += strconv.Itoa(v.RSSI) + ","
			rxRSSiList += strconv.Itoa(v.RxRSSI) + ","
			txNoAckList += strconv.Itoa(v.TxNoAck) + ","
			txTotalList += strconv.Itoa(v.TxTotal) + ","
		}
	}
	avgRSSi := tmp / cnt
	avgRxRSSi := tmp2 / cnt

	_, err = stmt1.Exec(timestamp, gwn, n0.ID, avgRSSi, avgRxRSSi, n0.Data.AppPER.LastSeq,
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
	defer stmt.Close()
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

	stmt, err := db.Prepare(`INSERT INTO NW_DATA_SET_LATENCY(TIMESTAMP, GATEWAY_NAME, SENSOR_ID, E2E_LATENCY) 
                VALUES(?,?,?,?)`)
	if err != nil {
		Error.Println(stmt, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(timestamp, gwn, n2.ID, n2.E2ELatency)
	if err != nil {
		Error.Println(stmt, err)
	}

	stmt2, err := db.Prepare(`UPDATE TOPOLOGY_DATA SET LAST_SEEN=? where (GATEWAY_NAME=? and SENSOR_ID=?) or SENSOR_ID=1`)
	if err != nil {
		Error.Println(stmt2, err)
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(timestamp, gwn, n2.ID)
	if err != nil {
		Error.Println(stmt2, err)
	}
}

// init logger and db
func init() {
	Info = log.New(infoHandle, "[*] INFO: ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "[!] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	db, _ = sql.Open("mysql", fmt.Sprintf("%v:%v@(db:3306)/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DB")))
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
	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
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
                        POWER VARCHAR(64) NOT NULL,
                        INDEX TOPO_INDEX(FIRST_APPEAR, GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table TOPOLOGY_DATA ready")

	// PARTITION_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS PARTITION_DATA (
                        TIMESTAMP BIGINT,
                        GATEWAY_NAME VARCHAR(16) NOT NULL,
                        ROWW SMALLINT NOT NULL,
                        TYPE VARCHAR(16) NOT NULL,
                        LAYER SMALLINT NOT NULL,
						CHANNEL_START SMALLINT NOT NULL,
                        CHANNEL_END SMALLINT NOT NULL,
                        START SMALLINT NOT NULL,
                        END SMALLINT NOT NULL,
                        INDEX PART_INDEX(TIMESTAMP, GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table PARTITION_DATA ready")

	// HARP_PARTITION_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS HARP_PARTITION_DATA (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		TYPE VARCHAR(16) NOT NULL,
		START SMALLINT NOT NULL,
		END SMALLINT NOT NULL,
		INDEX PART_INDEX(TIMESTAMP, GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table HARP_PARTITION_DATA ready")

	// HARP_SP_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS HARP_SP_DATA (
		TIMESTAMP BIGINT,
		GATEWAY_NAME VARCHAR(16) NOT NULL,
		SENSOR_ID SMALLINT UNSIGNED NOT NULL,
		LAYER SMALLINT NOT NULL,
		TS_START SMALLINT NOT NULL,
		TS_END SMALLINT NOT NULL,
		CH_START SMALLINT NOT NULL,
		CH_END SMALLINT NOT NULL,
		INDEX PART_INDEX(TIMESTAMP, GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table HARP_SP_DATA ready")

	// SCHEDULE_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS SCHEDULE_DATA (
                        TIMESTAMP BIGINT,
                        GATEWAY_NAME VARCHAR(16) NOT NULL,
                        SLOT_OFFSET SMALLINT NOT NULL,
                        CHANNEL_OFFSET SMALLINT NOT NULL,
                        SUBSLOT_OFFSET SMALLINT NOT NULL,
                        SUBSLOT_PERIOD SMALLINT NOT NULL,
                        TYPE VARCHAR(16) NOT NULL,
                        Layer SMALLINT,
                        SENDER SMALLINT,
                        RECEIVER INT,
                        IS_OPTIMAL SMALLINT,
                        INDEX SCHED_INDEX(TIMESTAMP,GATEWAY_NAME(16)));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table SCHEDUELE_DATA ready")

	// SENSOR_DATA
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS SENSOR_DATA (
                TIMESTAMP BIGINT,
                GATEWAY_NAME VARCHAR(16) NOT NULL,
                SENSOR_ID SMALLINT UNSIGNED NOT NULL,
                TEMP FLOAT NOT NULL,
				HUMIDITY FLOAT NOT NULL,
				ULTRASONIC FLOAT NOT NULL,
				LVDT1 FLOAT NOT NULL,
				LVDT2 FLOAT NOT NULL,
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
                ASN_STAMP1 INT NOT NULL,
                CHANNEL TINYINT UNSIGNED NOT NULL,
                BAT FLOAT NOT NULL,
                A2A_LATENCY FLOAT NOT NULL,
                LAST_UPLINK_LATENCY FLOAT NOT NULL,
                INDEX SENSOR_INDEX(TIMESTAMP, GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table SENSOR_DATA ready")

	// NW_DATA_SET_PER_UCONN or network_data_0
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_PER_UCONN (
                TIMESTAMP BIGINT,
                GATEWAY_NAME VARCHAR(16) NOT NULL,
                SENSOR_ID SMALLINT UNSIGNED NOT NULL,
                CHANNEL_INFO TEXT,
                AVG_RSSI SMALLINT NOT NULL,
                AVG_RXRSSI SMALLINT NOT NULL,
                APP_PER_SENT_LAST_SEQ INT NOT NULL,
                APP_PER_SENT INT NOT NULL,
                APP_PER_SENT_LOST INT NOT NULL,
                TX_FAIL INT NOT NULL,
                TX_NOACK INT NOT NULL,
                TX_TOTAL INT NOT NULL,
                RX_TOTAL INT NOT NULL,
                TX_LENGTH_TOTAL INT NOT NULL,
                MAC_TX_NOACK_DIFF INT NOT NULL,
                MAC_TX_TOTAL_DIFF INT NOT NULL,
                MAC_RX_TOTAL_DIFF INT NOT NULL,
                MAC_TX_LENGTH_TOTAL_DIFF INT NOT NULL,
                APP_PER_LOST_DIFF INT NOT NULL,
                APP_PER_SENT_DIFF INT NOT NULL,
                INDEX NW_PER_INDEX(TIMESTAMP, GATEWAY_NAME));`)
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
                TX_TOTAL VARCHAR(64) NOT NULL,
                INDEX CHINFO_INDEX(TIMESTAMP,GATEWAY_NAME));`)
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
                NUM_PARENT_CHANGE INT UNSIGNED NOT NULL,
                NUM_SYNC_LOST INT UNSIGNED NOT NULL,
                AVG_DRIFT INT UNSIGNED NOT NULL,
                MAX_DRIFT INT UNSIGNED NOT NULL,
                NUM_MAC_OUT_OF_BUFFER INT UNSIGNED NOT NULL,
                NUM_UIP_RX_LOST INT UNSIGNED NOT NULL,
                NUM_LOWPAN_TX_LOST INT UNSIGNED NOT NULL,
                NUM_LOWPAN_RX_LOST INT UNSIGNED NOT NULL,
                NUM_COAP_RX_LOST INT UNSIGNED NOT NULL,
                NUM_COAP_OBS_DIS INT UNSIGNED NOT NULL,
                INDEX NW_INFO_INDEX(TIMESTAMP,GATEWAY_NAME));`)
	if err != nil {
		Error.Panicln(err)
	}
	Info.Println("Table NW_DATA_SET_INFO ready")

	// NW_DATA_SET_LATENCY or network_data_2
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS NW_DATA_SET_LATENCY (
                TIMESTAMP BIGINT,
                GATEWAY_NAME VARCHAR(16) NOT NULL,
                SENSOR_ID SMALLINT UNSIGNED NOT NULL,
                E2E_LATENCY FLOAT NOT NULL);`)
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

	partition struct {
		Data []dataPartition `json:"data"`
	}
	dataPartition struct {
		Row      int    `json:"row"`
		Type     string `json:"type"`
		Layer    int    `json:"layer"`
		Channels [2]int `json:"channels"`
		Range    [2]int `json:"range"`
	}

	// HARP sub-partition
	subpartition struct {
		ID   int      `json:"_id"`
		Data []dataSP `json:"data"`
	}
	dataSP struct {
		Layer   int `json:"layer"`
		TsStart int `json:"ts_start"`
		TsEnd   int `json:"ts_end"`
		ChStart int `json:"ch_start"`
		ChEnd   int `json:"ch_end"`
	}

	schedule struct {
		Data []dataSchedule `json:"data"`
	}
	dataSchedule struct {
		Slot    [2]int `json:"slot"`
		Subslot [2]int `json:"subslot"`
		Cell    struct {
			Type     string `json:"type"`
			Layer    int    `json:"layer"`
			Sender   int    `json:"sender"`
			Receiver int    `json:"receiver"`
		}
		IsOptimal int `json:"is_optimal"`
	}

	sensor struct {
		ID   int        `json:"_id"`
		Data dataSensor `json:"data"`
	}
	dataSensor struct {
		Temp       float32 `json:"temp"`
		Humidity   float32 `json:"humidity"`
		Ultrasonic float32 `json:"ultrasonic"`
		LVDT1      float32 `json:"lvdt1"`
		LVDT2      float32 `json:"lvdt2"`

		Rhum              int     `json:"rhum"`
		Lux               int     `json:"lux"`
		Press             int     `json:"press"`
		Accelx            float32 `json:"accelx"`
		Accely            float32 `json:"accely"`
		Accelz            float32 `json:"type"`
		LED               int     `json:"led"`
		Eh                int     `json:"eh"`
		Eh1               int     `json:"eh1"`
		CC2650Active      int     `json:"cc2650_active"`
		CC2650Sleep       int     `json:"cc2650_sleep"`
		RFTx              float32 `json:"rf_tx"`
		RFRx              float32 `json:"rf_rx"`
		MSP432Active      int     `json:"msp432_active"`
		MSP432Sleep       int     `json:"msp432_sleep"`
		GPSEnActive       float32 `json:"gpsen_active"`
		GPSEnSleep        float32 `json:"gpsen_sleep"`
		Others            int     `json:"others"`
		Sequence          int     `json:"sequence"`
		ASNStamp1         int     `json:"asn_stamp1"`
		Channel           int     `json:"channel"`
		Bat               float32 `json:"bat"`
		A2ALatency        float32 `json:"a2a_latency"`
		LastUplinkLatency float32 `json:"last_uplink_latency"`
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
		ID         int          `json:"_id"`
		E2ELatency float32      `json:"e2e_latency"`
		Data       networkData2 `json:"data"`
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
