package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	coreLogger "github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	pool2 "github.com/dontkillmyvibee/sourcershubgo/internal/core/repository/postgres/pool"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/middleware"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/server"
	"github.com/dontkillmyvibee/sourcershubgo/internal/features/users/repository/postgres"
	"github.com/dontkillmyvibee/sourcershubgo/internal/features/users/service"
	"github.com/dontkillmyvibee/sourcershubgo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := coreLogger.NewLogger(coreLogger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := pool2.NewConnectionPool(ctx, pool2.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := postgres.NewUsersRepository(pool)
	usersService := service.NewUsersService(usersRepository)
	usersTransportHTTP := http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing http server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Panic(),
		middleware.Trace(),
	)
	apiVersionRouter := server.NewApiVersionRouter(server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error:", zap.Error(err))
	}
}
