package main

import (
	"net/url"
	"os"

	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"github.com/msales/transporter"
	"github.com/msales/transporter/kafka"
	"gopkg.in/inconshreveable/log15.v2"
)

// Application =============================

func newApplication(c *Context) (*transporter.Application, error) {
	consumer, err := newKafkaConsumer(c)
	if err != nil {
		return nil, err
	}

	app := transporter.NewApplication()
	app.Consumer = consumer

	return app, nil
}

// Consumer ================================

func newKafkaConsumer(c *Context) (*kafka.Consumer, error) {
	return kafka.New(
		c,
		kafka.WithBrokers(c.StringSlice(FlagKafkaBrokers)),
		kafka.WithGroupId(c.String(FlagKafkaGroupId)),
		kafka.WithTopics(c.StringSlice(FlagKafkaTopics)),
		kafka.WithBufferSize(c.Int(FlagBufferSize)),
	)
}

// Logger ==================================

func newLogger(c *Context) (log15.Logger, error) {
	lvl := c.String(FlagLogLevel)
	v, err := log15.LvlFromString(lvl)
	if err != nil {
		return nil, err
	}

	h := log15.LvlFilterHandler(v, log15.StreamHandler(os.Stderr, log15.LogfmtFormat()))
	if lvl == "debug" {
		h = log15.CallerFileHandler(h)
	}

	l := log15.New()
	l.SetHandler(log15.LazyHandler(h))

	return l, nil
}

// Stats ===================================

func newStats(c *Context) (stats.Stats, error) {
	dsn := c.String(FlagStats)
	if dsn == "" {
		return stats.Null, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "statsd":
		return newStatsdStats(uri.Host)

	case "l2met":
		return newL2metStats(c.logger), nil

	default:
		return stats.Null, nil
	}
}

func newStatsdStats(addr string) (stats.Stats, error) {
	return stats.NewStatsd(addr, "transporter")
}

func newL2metStats(log log.Logger) stats.Stats {
	return stats.NewL2met(log, "transporter")
}
