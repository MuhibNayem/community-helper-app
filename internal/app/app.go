package app

import (
	"fmt"

	"github.com/MuhibNayem/community-helper-app/internal/api"
	"github.com/MuhibNayem/community-helper-app/internal/api/handlers"
	"github.com/MuhibNayem/community-helper-app/internal/api/middleware"
	"github.com/MuhibNayem/community-helper-app/internal/config"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services/memory"
	"github.com/MuhibNayem/community-helper-app/internal/platform/server"
)

type App struct {
	cfg    *config.Config
	server *server.HTTPServer
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	store := memory.NewStore()

	handlerSet := api.HandlerSet{
		Auth:     handlers.NewAuthHandler(store),
		Users:    handlers.NewUsersHandler(store),
		Requests: handlers.NewRequestsHandler(store),
		Matches:  handlers.NewMatchesHandler(store),
		Health:   handlers.NewHealthHandler(cfg.Env),
	}

	authMiddleware := middleware.NewAuthMiddleware(store)

	router := api.NewRouter(cfg, handlerSet, authMiddleware)

	httpServer := server.NewHTTPServer(cfg.HTTPPort, router)

	return &App{
		cfg:    cfg,
		server: httpServer,
	}, nil
}

func (a *App) Run() error {
	return a.server.Start()
}
