# tesla_wall_connector_exporter

Prometheus exporter for Tesla Wall Connector (Gen3) metrics.

## Running

This exporter needs to be able to reach the wall connector on the same network. All that's required by the exporter is to point the exporter at a Wall Connector's address. For example, by passing the flag `-twc.address=192.168.1.123`. Using a domain name is recommended, in case the Wall Connector's IP changes.

The exporter is intended to map 1:1 to a Wall Connector. If you have multiple Wall Connectors, you'll need to run an exporter for each.

### Example docker-compose snippet

```yaml
tesla_wall_connector_exporter:
  image: ghcr.io/benclapp/tesla_wall_connector_exporter:latest
  container_name: tesla_wall_connector_exporter
  restart: always
  command:
  - -twc.address=teslawallconnector_abc123.home.local
  ports:
  - 9859:9859
```

### Command Flags

```
Usage of /tesla_wall_connector_exporter:
  -twc.address string
        [REQUIRED] The address of the Tesla Wall Connector.
  -web.listen-address string
        Address to listen on for HTTP requests. (default ":9859")
  -web.metrics-path string
        Path to expose metrics on. (default "/metrics")
```
