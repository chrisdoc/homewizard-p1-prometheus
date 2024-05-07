# homewizard-p1-prometheus
> A [Prometheus](https://prometheus.io/) exporter for the [Homewizard P1 energy monitor](https://www.homewizard.nl/homewizard-wi-fi-p1-meter)
## Running the exporter
To run the application you need to specify the IP address of the P1 meter in your network as an environment variable `HOST`
```
HOST={IP_ADDRESS} go run main.go
```

You can also define a `TICK` in second which specifies in which interval the exporter will fetch data from the P1 monitor. The default tick interval is 10seconds.

The exporter will listen by default on port 9898 which can be changed by defining an environment variable `PORT`.

### Docker
The exporter can also be run as a docker container like
```
docker run --rm -p 9898:9898 -e HOST={ID_ADDRESS} chriskies/homewizard-p1-prometheus:latest
```

## Exported Data
The exporter exposes the following data to Prometheus
```
# HELP active_power_l1_w he active power on L1 in W
# TYPE active_power_l1_w gauge
active_power_l1_w 0
# HELP active_power_l2_w he active power on L2 in W
# TYPE active_power_l2_w gauge
active_power_l2_w 140
# HELP active_power_l3_w The active power on L3 in W
# TYPE active_power_l3_w gauge
active_power_l3_w 8
# HELP active_power_w The active power in W
# TYPE active_power_w gauge
active_power_w 148
# HELP total_gas_m3 The total gas consumption in m3
# TYPE total_gas_m3 gauge
total_gas_m3 0
# HELP total_power_export_t1_kwh The total power export on T1 in kWh
# TYPE total_power_export_t1_kwh gauge
total_power_export_t1_kwh 0
# HELP total_power_export_t2_kwh The total power export on T2 in kWh
# TYPE total_power_export_t2_kwh gauge
total_power_export_t2_kwh 0
# HELP total_power_import_t1_kwh The total power import on T1 in kWh
# TYPE total_power_import_t1_kwh gauge
total_power_import_t1_kwh 2000.479
# HELP total_power_import_t2_kwh The total power import on T2 in kWh
# TYPE total_power_import_t2_kwh gauge
total_power_import_t2_kwh 0
# HELP wifi_strength Wifi strength in Db
# TYPE wifi_strength gauge
wifi_strength 36
```
