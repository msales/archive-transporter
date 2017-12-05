# The Transporter

[![Go Report Card](https://goreportcard.com/badge/github.com/msales/transporter)](https://goreportcard.com/report/github.com/msales/transporter)
[![Build Status](https://travis-ci.org/msales/transporter.svg?branch=master)](https://travis-ci.org/msales/transporter)
[![GitHub release](https://img.shields.io/github/release/msales/transporter.svg)](https://github.com/msales/transporter/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/msales/transporter/master/LICENSE)

## Configuration

Transporter can be configured with command line flags and environment variables. 

##### Command Line Flags

| Flag | Options | Multiple Allowed | Description | Environment Variable |
| ---- | ------- | ---------------- | ----------- | -------------------- |
| --log.level | debug, info, warn, error | No | The log level to use. | TRANSPORTER_LOG_LEVEL |
| --stats | | No | The stats dsn to connect to | TRANSPORTER_STATS |
| --stats.tags | | Yes | The tags to attach to all metrics | TRANSPORTER_STATS_TAGS |
| --port | | No | The address to bind to for the http server. | TRANSPORTER_PORT |
| --buffer.size | | No | The size of each topics buffer. | TRANSPORTER_BUFFER_SIZE |
| --kafka.brokers | | Yes | The kafka seed brokers connect to. Format: 'ip:port'. | TRANSPORTER_KAFKA_BROKERS |
| --kafka.group-id | | No | The kafka group id to subscribe to. | TRANSPORTER_KAFKA_GROUP_ID |
| --kafka.topics | | Yes | The kafka topics to subscribe to. | TRANSPORTER_KAFKA_TOPICS |

##### Multi value environment variables

When using environment variables where multiple values are allowed, the values should be comma seperated.
E.g. ```--kafka.topics=foo --kafka.topics=bar``` should become ```TRANSPORTER_KAFKA_TOPICS=foo,bar```.

## HTTP Endpoints

#### GET /health

Gets the current health status of Transporter. Returns a 200 status code if healthy, otherwise a 500 status code

#### GET /:topic

Gets the next item in the topic queue. A timeout of 100ms is used when waiting for items in the queue. After the timeout
a 204 status code is returned. If the topic does not exist, 404 status code is returned

## Contributors

We're supposed to tell you how to contribute to transporter here.  
Since this is github: You know the drill - open issues, fork, create PRs, ...

## License

MIT-License. As is. No warranties whatsoever. Mileage may vary. Batteries not included.
