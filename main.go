package main

import (
	"fmt"
	"net/http"

	fast "github.com/ugozlave/gofast"
	"github.com/ugozlave/gofast/faster"
)

func main() {
	app := fast.New(faster.NewAppConfig())

	// logger
	fast.Log(app, faster.NewFastLogger)

	// controllers
	fast.Add(app, faster.NewHealthController)
	fast.Add(app, NewMyController)

	// middleware
	fast.Use(app, faster.NewLogMiddleware)

	// services
	fast.Register[IService](app, NewMyService)

	// config
	fast.Cfg(app, NewMyConfig)

	app.Run()
}

type MyController struct {
	service IService
}

func NewMyController(ctx *fast.BuilderContext) *MyController {
	return &MyController{
		service: fast.Get[IService](ctx, fast.Scoped),
	}
}

func (c *MyController) Prefix() string {
	return "/my"
}

func (c *MyController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", c.Get)
	return mux
}

func (c *MyController) Get(w http.ResponseWriter, r *http.Request) {
	c.service.DoSomething()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type IService interface {
	DoSomething()
}

type MyService struct {
	logger fast.Logger
	config fast.ConfigProvider[MyConfig]
}

func NewMyService(ctx *fast.BuilderContext) *MyService {
	return &MyService{
		logger: fast.GetLogger[MyService](ctx, fast.Scoped),
		config: fast.MustGetConfig[MyConfig](ctx, fast.Singleton),
	}
}

func (s *MyService) DoSomething() {
	s.logger.Inf("Hello World")
	s.logger.Dbg(fmt.Sprintf("Config setting: %s", s.config.Value().Setting))
}

type MyConfig struct {
	Setting string `json:"Setting"`
}

func NewMyConfig(_ *fast.BuilderContext) *faster.Config[MyConfig] {
	var v MyConfig
	v.Setting = "default"
	return faster.NewConfig(v)
}
