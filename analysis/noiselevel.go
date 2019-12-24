package main

import (
	"database/sql"
	"fmt"
	"math"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	dbAddr := fmt.Sprintf("root:0311@tcp(127.0.0.1:3306)/6tisch")
	db, _ = sql.Open("mysql", dbAddr)
	fmt.Println("db connected")
}

type topoData struct {
	SensorID int
	Parent   int
}

type perData struct {
	SensorID             int
	AvgRSSI              float64
	AvgRxRSSI            float64
	MacRxTotalDiff       float64
	MacTxTotalDiff       float64
	MacTxNoACKDiff       float64
	MacTxLengthTotalDiff float64
}

func main() {
	var t topoData
	var per perData

	tList := make([]topoData, 0)
	perList := make([]perData, 0)

	rows0, err := db.Query("select SENSOR_ID,PARENT from TOPOLOGY_DATA where FIRST_APPEAR in (select MAX(FIRST_APPEAR) from TOPOLOGY_DATA group by SENSOR_ID);")
	if err != nil {
		fmt.Println(err)
	}
	for rows0.Next() {
		rows0.Scan(&t.SensorID, &t.Parent)
		tList = append(tList, t)
	}

	rows1, err := db.Query("SELECT SENSOR_ID,AVG(AVG_RSSI),AVG(AVG_RXRSSI),AVG(MAC_RX_TOTAL_DIFF),AVG(MAC_TX_NOACK_DIFF),AVG(MAC_TX_TOTAL_DIFF),AVG(MAC_TX_LENGTH_TOTAL_DIFF) FROM NW_DATA_SET_PER_UCONN GROUP BY SENSOR_ID")
	if err != nil {
		fmt.Println(err)
	}
	for rows1.Next() {
		rows1.Scan(&per.SensorID, &per.AvgRSSI, &per.AvgRxRSSI, &per.MacRxTotalDiff, &per.MacTxNoACKDiff, &per.MacTxTotalDiff, &per.MacTxLengthTotalDiff)
		perList = append(perList, per)

	}
	for _, p := range perList {
		acRssi := float64(p.AvgRSSI)
		acTx := float64(p.MacRxTotalDiff)
		acLost := float64(p.MacTxNoACKDiff - p.MacTxTotalDiff - p.MacRxTotalDiff)

		var (
			txRssi   float64
			txTx     float64
			txLost   float64
			txLength float64
		)

		var parent = 0
		for _, t := range tList {
			if p.SensorID == t.SensorID {
				parent = t.Parent
				break
			}
		}
		for _, pp := range perList {
			if pp.SensorID == parent {
				txRssi = pp.AvgRxRSSI
				txTx = pp.MacTxTotalDiff
				txLost = pp.MacTxTotalDiff - pp.MacRxTotalDiff
				txLength = pp.MacTxLengthTotalDiff
			}
		}
		var noiseLevel float64
		if acTx > 0 && txTx > 0 {
			acNosie := math.Pow(10, (noiseCompute(acTx, acLost, acRssi, 20) / 10))
			txNoise := math.Pow(10, (noiseCompute(txTx, txLost, txRssi, txLength/txTx) / 10))
			noiseLevel = 10 * math.Log10(txNoise*txTx/(txTx+acTx)+acNosie*acTx/(txTx+acTx))
		} else if acTx > 0 {
			noiseLevel = noiseCompute(acTx, acLost, acRssi, 20)
		} else if txTx > 0 {
			noiseLevel = math.Pow(10, (noiseCompute(txTx, txLost, txRssi, txLength/txTx) / 10))
		} else {
			noiseLevel = -99.0
		}
		fmt.Println(acRssi, p.AvgRxRSSI, math.Round(noiseLevel))
	}
}

func noiseCompute(txTotal, lostTotal, rssi, length float64) float64 {
	return rssi - snrDb(lostTotal/txTotal, length)
}

func snrDb(plrIn, length float64) float64 {
	midSNR := 0.000000
	maxSNR := 4.000000
	minSNR := -4.000000
	midPLR := plr(midSNR, length)

	for math.Abs(minSNR-maxSNR) > 0.00001 {
		midSNR = (maxSNR + minSNR) / 2
		midPLR = plr(midPLR, length)
		if math.Abs(plrIn-midPLR) < 0.00001 {
			return midSNR
		} else if plrIn > midPLR {
			maxSNR = midSNR - 0.000001
		} else if plrIn < midPLR {
			minSNR = midSNR + 0.000001
		}
	}
	return midSNR
}

func plr(snrIn, length float64) float64 {
	bitErrorRate := 0.5 * math.Erfc(math.Sqrt(math.Pow(10, (snrIn/10)))) * 1.45
	para := 0.0

	for i := 1; i < 33; i++ {
		para = para + fact(32)/(fact(i)*fact(32-i))*math.Pow(bitErrorRate, float64(i))*(math.Pow(1-bitErrorRate, float64(32-i)))*p(i)
	}
	return 1 - math.Pow((1-para), 2*length)
}

func fact(n int) float64 {
	res := 1
	for i := 2; i < n+1; i++ {
		res *= i
	}
	return float64(res)
}

func p(n int) float64 {
	if n <= 5 {
		return 0.000000
	} else if n == 6 {
		return 0.002000
	} else if n == 7 {
		return 0.013400
	} else if n == 8 {
		return 0.052300
	} else if n == 9 {
		return 0.149800
	} else if n == 10 {
		return 0.347900
	} else if n == 11 {
		return 0.649600
	} else if n == 12 {
		return 0.915600
	} else if n == 13 {
		return 0.996800
	} else {
		return 1.000000
	}
}
