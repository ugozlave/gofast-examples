package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ugozlave/gofast/faster"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := faster.New()
	app.Run(ctx)

}
