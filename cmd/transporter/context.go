package main

import (
	"context"

	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"gopkg.in/urfave/cli.v1"
)

type netCtx context.Context

// Context provides instance of services the CLI consumes.
type Context struct {
	*cli.Context
	netCtx

	logger log.Logger
	stats  stats.Stats
}

func newContext(c *cli.Context) (ctx *Context, err error) {
	ctx = &Context{
		Context: c,
		netCtx:  context.Background(),
	}

	ctx.logger, err = newLogger(ctx)
	if err != nil {
		return
	}

	ctx.stats, err = newStats(ctx)
	if err != nil {
		return
	}

	if ctx.logger != nil {
		ctx.netCtx = log.WithLogger(ctx.netCtx, ctx.logger)
	}

	if ctx.stats != nil {
		ctx.netCtx = stats.WithStats(ctx.netCtx, ctx.stats)
	}

	return
}
