# github.com/sesopenko/ecowitt-proxy

A Go-based proxy server for forwarding Ecowitt GW1200 weather station data to multiple endpoints, including
Home Assistant and Hubitat, with customizable configuration and logging support.

This Go-based proxy server allows you to overcome the limitation of Ecowitt gateways, which support only a
single custom data endpoint. By using this app, you can simultaneously send your weather data to multiple
destinations, such as Hubitat, Home Assistant, and an InfluxDB ingestion app. This ensures that all your smart
home systems and data analytics tools receive the weather data they need without any conflicts or manual
intervention.

## Setup

1. Copy `config.example.yml` to `config.yml` and enter your own details.  `config.yml` must be in your working directory.

# Licensed GNU GPL V3

This software is under the GNU GENERAL PUBLIC LICENSE VERSION 3. The license can be read in [LICENSE](LICENSE.txt). If
the license isn't included you may read it at
[https://www.gnu.org/licenses/gpl-3.0.txt](https://www.gnu.org/licenses/gpl-3.0.txt)

## Copyright

Copyright © Sean Esopenko 2024 All Rights Reserved