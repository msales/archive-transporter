package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

import _ "github.com/joho/godotenv/autoload"

// Flag constants declared for CLI use.
const (
	FlagPort     = "port"
	FlagLogLevel = "log.level"
	FlagStats    = "stats"

	FlagBufferSize = "buffer.size"

	FlagKafkaBrokers = "kafka.brokers"
	FlagKafkaGroupID = "kafka.group-id"
	FlagKafkaTopics  = "kafka.topics"
)

var commonFlags = []cli.Flag{
	cli.StringFlag{
		Name:   FlagLogLevel,
		Value:  "info",
		Usage:  "Specify the log level. You can use this to enable debug logs by specifying `debug`.",
		EnvVar: "TRANSPORTER_LOG_LEVEL",
	},
	cli.StringFlag{
		Name:   FlagStats,
		Value:  "",
		Usage:  "The stats backend to use. (e.g. statsd://localhost:8125)",
		EnvVar: "TRANSPORTER_STATS",
	},
}

var commands = []cli.Command{
	{
		Name:  "server",
		Usage: "Run the transporter",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:   FlagPort,
				Value:  "80",
				Usage:  "Specify the port to run the server on",
				EnvVar: "TRANSPORTER_PORT",
			},
			cli.IntFlag{
				Name:   FlagBufferSize,
				Value:  1000,
				Usage:  "The number of messages to buffer for each topic (default: 10000)",
				EnvVar: "TRANSPORTER_BUFFER_SIZE",
			},
			cli.StringSliceFlag{
				Name:   FlagKafkaBrokers,
				Usage:  "Specify the Kafka seed brokers",
				EnvVar: "TRANSPORTER_KAFKA_BROKERS",
			},
			cli.StringFlag{
				Name:   FlagKafkaGroupID,
				Value:  "transporter",
				Usage:  "Specify the Kafka consumer group id",
				EnvVar: "TRANSPORTER_KAFKA_GROUP_ID",
			},
			cli.StringSliceFlag{
				Name:   FlagKafkaTopics,
				Usage:  "Specify the Kafka topics to consume",
				EnvVar: "TRANSPORTER_KAFKA_TOPICS",
			},
		}, commonFlags...),
		Action: runServer,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "transporter"
	app.Usage = "A Kafka consumption HTTP server."
	app.Version = Version
	app.Commands = commands

	app.Run(os.Args)
}
