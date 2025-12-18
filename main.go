package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"slices"

	fast "github.com/ugozlave/gofast"
	"github.com/ugozlave/gofast/faster"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := faster.New()
	defer app.Shutdown()

	// config
	fast.Cfg(app, NewUserConfig)

	// controllers
	fast.Add(app, NewUserController)

	// middlewares
	fast.Use(app, NewUserMiddleware)

	// services
	fast.Register[IUserService](app, NewUserService)

	app.Run(ctx)
}

// config

type UserConfig struct {
	Users []string `json:"Users"`
}

func NewUserConfig(_ *fast.BuilderContext) *faster.FastConfig[UserConfig] {
	var v UserConfig
	v.Users = []string{"root:root"}
	return faster.NewConfig(v, "UserSettings")
}

// controllers

type UserController struct {
	logger  fast.Logger
	service IUserService
}

func NewUserController(ctx *fast.BuilderContext) *UserController {
	return &UserController{
		logger:  fast.MustGetLogger[UserController](ctx, fast.Scoped),
		service: fast.MustGet[IUserService](ctx, fast.Scoped),
	}
}

func (c *UserController) Prefix() string {
	return "/my"
}

func (c *UserController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", c.Get)
	return mux
}

func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
	username := c.service.Get("username")
	c.logger.Inf("authorized: " + username.(string))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// middlewares

type UserMiddleware struct {
	config  fast.Config[UserConfig]
	service IUserService
}

func NewUserMiddleware(ctx *fast.BuilderContext) *UserMiddleware {
	return &UserMiddleware{
		config:  fast.MustGetConfig[UserConfig](ctx, fast.Singleton),
		service: fast.MustGet[IUserService](ctx, fast.Scoped),
	}
}

func (m *UserMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := m.config.Value().Users
		username, password, _ := r.BasicAuth()
		if !slices.Contains(users, username+":"+password) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		m.service.Set("username", username)
		next.ServeHTTP(w, r)
	})
}

// services

type IUserService interface {
	Get(string) any
	Set(string, any)
}

type UserService struct {
	data map[string]any
}

func NewUserService(ctx *fast.BuilderContext) *UserService {
	return &UserService{
		data: make(map[string]any),
	}
}

func (s *UserService) Get(key string) any {
	return s.data[key]
}

func (s *UserService) Set(key string, value any) {
	s.data[key] = value
}
