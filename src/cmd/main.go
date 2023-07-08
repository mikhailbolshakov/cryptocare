package main

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/bootstrap"
	kitContext "github.com/mikhailbolshakov/cryptocare/src/kit/context"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"os"
	"os/signal"
	"syscall"
)

//@title CryptoCare API
//@version 1.0
//@description CryptoCare service allows improving your arbitrage trading experience
//@contact.name Api service support
//@contact.email support@cryptocare.io
//
//@BasePath /api
func main() {

	// init context
	ctx := kitContext.NewRequestCtx().Empty().WithNewRequestId().ToContext(context.Background())

	// create a new service
	s := bootstrap.New()

	l := service.L().Mth("main").Inf("created")

	// init service
	if err := s.Init(ctx); err != nil {
		l.E(err).St().Err("initialization")
		os.Exit(1)
	}

	l.Inf("initialized")

	// start listening
	if err := s.Start(ctx); err != nil {
		l.E(err).St().Err("listen")
		os.Exit(1)
	}

	l.Inf("listening")

	// handle app close
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Inf("graceful shutdown")
	s.Close(ctx)
	os.Exit(0)
}
