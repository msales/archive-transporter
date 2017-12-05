package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/msales/transporter"
	"github.com/msales/transporter/server"
	"github.com/msales/transporter/server/middleware"
	"gopkg.in/urfave/cli.v1"
)

func runServer(c *cli.Context) {
	ctx, err := newContext(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer app.Close()

	port := c.String(FlagPort)
	s := newServer(ctx, app)
	ctx.logger.Info(fmt.Sprintf("Starting server on port %s", port))
	log.Fatal(http.ListenAndServe(":"+port, s))
}

func newServer(ctx *Context, app *transporter.Application) http.Handler {
	s := server.New(app)

	h := middleware.Common(s)
	return middleware.WithContext(ctx, h)
}
