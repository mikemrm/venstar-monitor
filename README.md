# Venstar Monitor

A monitoring solution for Venstar.

## Support outputs

One or more of the following output formats may be defined, however only one of
each can be used.

* `json`: Dumps the results from querying the device to stdout.
* `influx`: Writes the results to the specified influxdb host.
* `prometheus`: Starts an http server so prometheus can scrape for metrics.

Both `json` and `influx` output methods use interval monitoring while
`prometheus` does pull based scraping, so there can be inconsistencies in
values when comparing as it's unlikely the scrapes from the two approaches
occur at the same time.

## Usage

```shell
$ ./venstar-monitor --help
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
      --prometheus.listen string    Prometheus exporter listen address (default ":9872")
  -t, --thermostat strings          One or more thermostat hosts to monitor
```

## TODOs:

- Add proper logging rather than simply printing
- Add tests
- Setup ci
- Add dockerfile
- Add sensors, runtimes and alerts to influx/prometheus outputs

## JSON Output

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

## InfluxDB Output

```shell
$ ./venstar-monitor -t 192.168.1.102 -o influx \
    --influx.addr 'http://influx:8086' \
    --influx.database homestats \
    --influx.measurement venstar_thermostat
```

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

## Prometheus Output

```shell
$ ./venstar-monitor -t 192.168.1.102 -o prometheus
Starting prometheus server...
Servers started!
```

```shell
$ curl -s http://venstar-monitor:9872/metrics
# HELP venstar_thermostat_active_stage Thermostat current heat/cool active stage
# TYPE venstar_thermostat_active_stage gauge
venstar_thermostat_active_stage{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_available_modes Thermostat available modes 0:all 1:heat/cool 2:heat 3:cool
# TYPE venstar_thermostat_available_modes gauge
venstar_thermostat_available_modes{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_away Thermostat away 0:home 1:away
# TYPE venstar_thermostat_away gauge
venstar_thermostat_away{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_cool_temp Thermostat set cool temperature
# TYPE venstar_thermostat_cool_temp gauge
venstar_thermostat_cool_temp{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 74
# HELP venstar_thermostat_cool_temp_max Thermostat set cool temperature max
# TYPE venstar_thermostat_cool_temp_max gauge
venstar_thermostat_cool_temp_max{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 99
# HELP venstar_thermostat_cool_temp_min Thermostat set cool temperature min
# TYPE venstar_thermostat_cool_temp_min gauge
venstar_thermostat_cool_temp_min{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 35
# HELP venstar_thermostat_dehumidify_setpoint Thermostat dehumidify percentage setpoint
# TYPE venstar_thermostat_dehumidify_setpoint gauge
venstar_thermostat_dehumidify_setpoint{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 99
# HELP venstar_thermostat_fan Thermostat fan setting 0:auto 1:on
# TYPE venstar_thermostat_fan gauge
venstar_thermostat_fan{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 1
# HELP venstar_thermostat_fan_state Thermostat current fan state 0:off 1:on
# TYPE venstar_thermostat_fan_state gauge
venstar_thermostat_fan_state{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_force_unoccupied Thermostat force unoccupied 0:off 1:on
# TYPE venstar_thermostat_force_unoccupied gauge
venstar_thermostat_force_unoccupied{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_heat_temp Thermostat set heat temperature
# TYPE venstar_thermostat_heat_temp gauge
venstar_thermostat_heat_temp{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 70
# HELP venstar_thermostat_heat_temp_max Thermostat set heat temperature max
# TYPE venstar_thermostat_heat_temp_max gauge
venstar_thermostat_heat_temp_max{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 99
# HELP venstar_thermostat_heat_temp_min Thermostat set heat temperature min
# TYPE venstar_thermostat_heat_temp_min gauge
venstar_thermostat_heat_temp_min{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 35
# HELP venstar_thermostat_holiday Thermostat holiday 0:not observing 1:observing
# TYPE venstar_thermostat_holiday gauge
venstar_thermostat_holiday{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_humidity Thermostat humidity
# TYPE venstar_thermostat_humidity gauge
venstar_thermostat_humidity{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 49
# HELP venstar_thermostat_humidity_enabled Thermostat humidity enabled 0:disabled 1:enabled
# TYPE venstar_thermostat_humidity_enabled gauge
venstar_thermostat_humidity_enabled{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_humidity_setpoint Thermostat humidity percentage setpoint
# TYPE venstar_thermostat_humidity_setpoint gauge
venstar_thermostat_humidity_setpoint{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_mode Thermostat mode 0:off 1:heat 2:cool 3:auto
# TYPE venstar_thermostat_mode gauge
venstar_thermostat_mode{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 3
# HELP venstar_thermostat_override Thermostat override 0:off 1:on
# TYPE venstar_thermostat_override gauge
venstar_thermostat_override{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_override_remaining_seconds Thermostat override remaining seconds
# TYPE venstar_thermostat_override_remaining_seconds gauge
venstar_thermostat_override_remaining_seconds{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_schedule Thermostat schedule 0:inactive 1:active
# TYPE venstar_thermostat_schedule gauge
venstar_thermostat_schedule{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_schedule_part Thermostat schedule part 0:morning 1:day 2:evening 3:night 255:inactive
# TYPE venstar_thermostat_schedule_part gauge
venstar_thermostat_schedule_part{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 255
# HELP venstar_thermostat_setpoint_delta Thermostat setpoint delta
# TYPE venstar_thermostat_setpoint_delta gauge
venstar_thermostat_setpoint_delta{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 4
# HELP venstar_thermostat_space_temp Thermostat current space temperature
# TYPE venstar_thermostat_space_temp gauge
venstar_thermostat_space_temp{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 74
# HELP venstar_thermostat_state Thermostat current state 0:idle 1:heating 2:cooling 3:lockout 4:error
# TYPE venstar_thermostat_state gauge
venstar_thermostat_state{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
# HELP venstar_thermostat_temp_units Thermostat temperature units 0:fahrenheit 1:celsius
# TYPE venstar_thermostat_temp_units gauge
venstar_thermostat_temp_units{host="192.168.1.102",model="COLORTOUCH",name="Thermostat",type="commercial"} 0
```
