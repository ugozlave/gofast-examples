package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/ugozlave/gofast"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app, _ := gofast.New()
	// Add a controller
	gofast.Add(app, MyControllerBuilder())
	// Add a config
	gofast.Cfg(app, gofast.ConfigBuilder(MyConfig{Words: []string{"my", "default", "values"}}))
	// Add a service
	gofast.Register[IService](app, MyServiceBuilder())

	app.Run(ctx)

}

type MyController struct {
	service IService
}

func MyControllerBuilder() gofast.Builder[*MyController] {
	return func(ctx *gofast.BuilderContext) *MyController {
		return &MyController{
			service: gofast.MustGet[IService](ctx, gofast.Scoped),
		}
	}
}

func (c *MyController) Prefix() string {
	return "my"
}

func (c *MyController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", c.handle)
	return mux
}

func (c *MyController) handle(w http.ResponseWriter, r *http.Request) {
	c.service.DoSomething()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type IService interface {
	DoSomething()
}

type MyService struct {
	config gofast.Config[MyConfig]
	logger gofast.Logger
}

func MyServiceBuilder() gofast.Builder[*MyService] {
	return func(ctx *gofast.BuilderContext) *MyService {
		return &MyService{
			config: gofast.MustGetConfig[MyConfig](ctx, gofast.Singleton),
			logger: gofast.MustGetLogger[MyService](ctx, gofast.Scoped),
		}
	}
}

func (s *MyService) DoSomething() {
	s.logger.Inf(fmt.Sprintf("Doing something with config: %v", s.config.Value().Words))
}

type MyConfig struct {
	Words []string `json:"Words"`
}

func (c MyConfig) Path() []string {
	return []string{"CustomConfig", "MyConfig"}
}
