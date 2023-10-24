# Enphase Envoy Golang Client

This is a prometheus collector for pulling metrics from an Envoy Enphase unit. The collector
utilizes the local interface exposed by the device rather than the Enlighten API. Enphase units
are embedded devices, so the collector is implemented as a proxy collector similar to the
snmp_exporter tool.

<https://enphase.com/en-us/support/what-envoy>

## Example

```yml
  - job_name: 'prometheus-envoy'
    static_configs:
      - targets:
        - '192.168.1.40'
        - '192.168.1.41'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:2112  # The prometheus-smarthome's real hostname:port.
```

## Building and running

```sh
cd cmd/prometheus-envoy
go build
./prometheus-envoy -port 2112 -token <jwt>
```

## Installing

```sh
make && sudo make install-svc
mkdir -p /etc/prometheus-envoy/
# fetch a JWT from the envoy
echo "<jwt>" | sudo tee /etc/prometheus/token
sudo systemctl enable prometheus-envoy.service
sudo systemctl start prometheus-envoy.service
```

## Listening on all interfaces

By default, this will only listen on the loopback interface. This can be problematic
if prometheus is running on a different host. To listen on all interfaces, add `-listen=0.0.0.0`
to the command line argument. This is done by default in the systemd service definition.

## License

This library is provided under the [MIT License](LICENSE.md)
