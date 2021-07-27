# 6TiSCH Testbed

## JSON from Gateway to Cloud

### Topology Data

New device topology data

```json
msg": {
    "data": {
      "_id": 4,
      "address": "2001:db8:1234:ffff::ff:fe00:4",
      "parent": 67,
      "eui64": "00-12-4b-00-12-05-21-e0",
      "gps": [
        41.80670888614393,
        -72.25271098086125
      ],
      "type": "SensorTag",
      "power": "battery"
    }
}

```

### Network basic data of device

```json
"msg": {
    "_id": 37,
    "data": {
      "curParent": 7,
      "numParentChange": 2,
      "numSyncLost": 2,
      "avgDrift": 157,
      "maxDrift": 396,
      "numMacOutOfBuffer": 0,
      "numUipRxLost": 0,
      "numLowpanTxLost": 0,
      "numLowpanRxLost": 0,
      "numCoapRxLost": 0,
      "numCoapObsDis": 0
    }
  }
```

### Network latency data of device

```json
"msg": {
    "_id": 37,
    "rtt": 3.924
  }
```

### Network advanced data of device

64 channel information

```json
"msg":{ 
   "_id":3,
   "data":{ 
      "ch":{ 
         "0":{ 
            "rssi":1,
            "rxRssi":2,
            "txNoAck":3,
            "txTotal":0
         },
         ...
      },
      "appPer":{ 
         "last_seq":160,
         "sent":27,
         "lost":0
      },
      "txFail":116,
      "txNoAck":2,
      "txTotal":89,
      "rxTotal":87,
      "txLengthTotal":4792,
      "macTxNoAckDiff":2,
      "macTxTotalDiff":36,
      "macRxTotalDiff":34,
      "macTxLengthTotalDiff":2087,
      "appLostDiff":0,
      "appSentDiff":12
   }
}
}'
```

### Sensing data

```json
"msg": {
    "_id": 18,
    "data": {
      "temp": 21.5,
      "rhum": 0,
      "lux": 55,
      "press": 972,
      "accelx": -0.66,
      "accely": -0.63,
      "accelz": -0.24,
      "led": 0,
      "eh": 0,
      "eh1": 0,
      "cc2650_active": 2,
      "cc2650_sleep": 2,
      "rf_tx": 11.5,
      "rf_rx": 30.36,
      "msp432_active": 0,
      "msp432_sleep": 0,
      "gpsen_active": 16.5,
      "gpsen_sleep": 49.8,
      "others": 5,
      "sequence": 145,
      "asn_stamp1": 0,
      "asn_stamp2": 0,
      "channel": 14,
      "bat": 2.69,
      "latency": 1.05
    }
  }
```

## JSON retrieved for Frontend (Cloud API)

### Common JSON Format

```text
{
    // results
    "data": {},
    // 1: success, 0: some error occured
    "flag": 1,
    // "success" or error message
    "msg": "success"
}
```
### Common Parameters

Time range: `range=hour|day|week|month`, default query all data.

Advanced info flag: `advanced=1|0`

### GET /api/gateway

Get gateway name list

Params: `range`

Example:

```text
data": [
    "UCONN_GW",
    "UCONN_GWX"
],
```

### GET /api/:gateway/topology

Get topology data for a gateway.

Params: `range`

Example:

```text
"data": [
    ...
    {
        "datetime": "2019-11-06 16:55:34",
        "timestamp": 1573077334467,
        "gateway": "UCONN_GW",
        "sensor_id": 1,
        "address": "2001:db8:1234:ffff::ff:fe00:1",
        "parent": 0,
        "eui64": "",
        "position": {
            "lat": 41.806737363265306,
            "lng": -72.25269334242424
        },
        "type": "LaunchPad",
        "power": "USB"
    },
    ...
],
```

### GET /api/:gateway/nwstat

Get network statistic data for all sensors under a gateway.

Params: `range`

Example:

```text
"data": [
    ...
    {
        "sensor_id": 3,
        "gateway": "UCONN_GW",
        "avg_rtt": 1.4585000276565552,
        "avg_mac_tx_total_diff": 27.4,
        "avg_mac_tx_noack_diff": 0,
        "avg_app_per_sent_diff": 8.8,
        "avg_app_per_lost_diff": 0
    },
    ...
],
```

### GET /api/:gateway/nwstat/:sensorID

Get detail network statistic data for one sensor under a gateway.

Params: `range, advanced`

Example:

advanced=0

```text
"data": [
    ...
    {
        "timestamp": 1573077385010,
        "avg_rssi": -44,
    },
    ...
]
```

advanced=1

```text
"data": [
    ...
    {
        "timestamp": 1573077385010,
        "mac_tx_total_diff": 35,
        "mac_tx_noack_diff": 1,
        "app_per_sent_diff": 12,
        "app_per_lost_diff": 0
    }
    ...
]
```

### GET /api/:gateway/battery

Get battery data for all sensors under a gateway.

Params: none

Example:

```text
"data": [
    ...
    {
        "gateway": "UCONN_GW",
        "sensor_id": 3,
        "avg_cc2650_active": 2,
        "avg_cc2650_sleep": 2,
        "avg_rf_rx": 90.34433875893647,
        "avg_rf_tx": 13.936226430928931,
        "bat_remain": "96%:31D"
    },
    ...
]
```