# Web Backend of 6TiSCH Cloud

## RESTful API

### Common Rules

#### Common JSON Format

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
#### Common Parameters

Time range: `range=hour|day|week|month`, default query all data.

Advanced info flag: `advanced=1|0`

### APIs

#### GET /api/gateway

Get gateway name list

Params: `range`

Example:

```text
data": [
    "UCONN_GW",
    "UCONN_GWX"
],
```

#### GET /api/:gateway/topology

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

#### GET /api/:gateway/nwstat

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

#### GET /api/:gateway/nwstat/:sensorID

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

#### GET /api/:gateway/battery

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