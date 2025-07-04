package main

import (
	"net/http"

	fast "github.com/ugozlave/gofast"
	"github.com/ugozlave/gofast/faster"
)

func main() {
	app := fast.New()

	// controllers
	fast.Add(app, faster.NewHealthController)
	fast.Add(app, NewMyController)

	// middleware
	fast.Use(app, faster.NewLogMiddleware)

	// services
	fast.Register[IService](app, NewMyService)

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
}

func NewMyService(ctx *fast.BuilderContext) *MyService {
	return &MyService{
		logger: fast.GetLogger[MyService](ctx, fast.Scoped),
	}
}

func (s *MyService) DoSomething() {
	s.logger.Inf("Hello World")
}
