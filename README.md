# Venstar Monitor

A monitoring solution for Venstar.

## Support outputs

* `json`: Dumps the results from querying the device to stdout.
* `influx`: Writes the results to the specified influxdb host.

## Usage

```shell
$ venstar-monitor --help
A service to monitor your Venstar devices

Usage:
  venstar-monitor [flags]

Flags:
  -h, --help                        help for venstar-monitor
      --influx.addr string          InfluxDB address (default "http://localhost:8086")
      --influx.database string      InfluxDB database (required)
      --influx.measurement string   InfluxDB measurement (required)
      --influx.pass string          InfluxDB authentication password
      --influx.retention string     InfluxDB retention policy (default "autogen")
      --influx.user string          InfluxDB authentication username
  -i, --interval duration           How frequent to poll devices for updates (default 5s)
  -o, --output strings              One or more output formats (default [json])
  -t, --thermostat strings          One or more thermostat hosts to monitor
```

## JSON Example

```shell
$ venstar-monitor -t 192.168.1.102 | jq -r '[
  .Device.Address, .Timestamp, .Thermostat.QueryInfo.spacetemp,
  ["idle", "heating", "cooling"][.Thermostat.QueryInfo.state]
] | @tsv'
192.168.1.102	2020-09-19T19:42:23.520772152Z	75	cooling
192.168.1.102	2020-09-19T19:42:28.520761588Z	75	cooling
192.168.1.102	2020-09-19T19:42:33.520777056Z	74	idle
192.168.1.102	2020-09-19T19:42:38.520766028Z	74	idle
```

## InfluxDB Example

```shell
> select * from thermostat
name: thermostat
time                cool_temp fan fan_state heat_temp humidity name       name_1     space_temp state system_model system_type temp_units
----                --------- --- --------- --------- -------- ----       ------     ---------- ----- ------------ ----------- ----------
1600544695000000000 74        on  off       70        48       Thermostat Thermostat 74         idle  COLORTOUCH   commercial  fahrenheit
1600544700000000000 74        on  off       70        48       Thermostat Thermostat 74         idle  COLORTOUCH   commercial  fahrenheit
1600544705000000000 74        on  off       70        48       Thermostat Thermostat 74         idle  COLORTOUCH   commercial  fahrenheit
1600544710000000000 74        on  off       70        48       Thermostat Thermostat 74         idle  COLORTOUCH   commercial  fahrenheit
```