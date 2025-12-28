package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ugozlave/gofast"
)

type DummyConfig[T any] struct {
	value T
}

func NewConfig[T any](v T) *DummyConfig[T] {
	return &DummyConfig[T]{value: v}
}

func (c *DummyConfig[T]) Value() T {
	return c.value
}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := gofast.AppConfig{}
	cfg.Name = "basic"
	cfg.Server.Port = 8080

	app := gofast.New(NewConfig(cfg))

	app.Run(ctx)

}
