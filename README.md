# github.com/sesopenko/ecowitt-proxy

A Go-based proxy server for forwarding Ecowitt GW1200 weather station data to multiple endpoints, including
Home Assistant and Hubitat, with customizable configuration and logging support.

This Go-based proxy server allows you to overcome the limitation of Ecowitt gateways, which support only a
single custom data endpoint. By using this app, you can simultaneously send your weather data to multiple
destinations, such as Hubitat, Home Assistant, and an InfluxDB ingestion app. This ensures that all your smart
home systems and data analytics tools receive the weather data they need without any conflicts or manual
intervention.

## Setup

### Configuration

Example `config.yml`
```yaml
targets:
  # An example of forwarding to home assistant https://www.home-assistant.io/integrations/ecowitt
  - name: home-assistant
    host_addr: http://192.168.1.20:8123/api/webhook/12345890jklfd89043jkl
  # an example of forwarding to https://github.com/sesopenko/ecowitt-to-influxdb
  - name: ecowitt-to-influx
    host_addr: http://192.168.1.22:20555/data/report/
  # the following is untested for hubitat at this point because I don't have one yet.
  # based off docs from https://github.com/padus/ecowitt
  - name: hubitat
    host_addr: http://192.168.1.21:39501/data
  # an example of forwarding to ecowitt2mqtt https://github.com/bachya/ecowitt2mqtt
  - name: ecowitt2mqtt
    host_addr: http://192.168.1.23:8080/data/report/

server:
  path: /api/webhook/someurl
  verbose: false
  # Set this to true to skip tls verification when sending data to targets.
  # Don't set this to true in production.
  tls_insecure_skip_verify: false
```

Create `config.yml` and enter your own details.  `config.yml` must be in your working directory.

## Running from docker

```bash
docker run -d -p 8123:8123 --name ecowitt-proxy -v ./config.yml:/app/config.yml:ro sesopenko/ecowitt-proxy
```

## Running using docker-compose

`docker-compose.yml`
```
version: '3.8'

services:
  ecowitt-proxy:
    image: sesopenko/ecowitt-proxy
    restart: unless-stopped
    ports:
      - "8123:8123"
    volumes:
      - ./config.yml:/app/config.yml
```

# Licensed GNU GPL V3

This software is under the GNU GENERAL PUBLIC LICENSE VERSION 3. The license can be read in [LICENSE](https://github.com/sesopenko/ecowitt-proxy/blob/main/LICENSE.txt). If
the license isn't included you may read it at
[https://www.gnu.org/licenses/gpl-3.0.txt](https://www.gnu.org/licenses/gpl-3.0.txt)

## Copyright

Copyright © Sean Esopenko 2024 All Rights Reserved