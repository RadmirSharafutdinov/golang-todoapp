package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	core_logger "study/internal/core/logger"
	core_postgres_poll "study/internal/core/repository/postgres/pool"
	core_http_middleware "study/internal/core/transport/http/middleware"
	core_http_server "study/internal/core/transport/http/server"
	users_postgres_repository "study/internal/features/users/repository/postgres"
	users_service "study/internal/features/users/service"

	users_transport_http "study/internal/features/users/transport/http"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()
	logger.Debug("initiazling postgres connection pool")

	pool, err := core_postgres_poll.NewConnectionPool(
		ctx,
		core_postgres_poll.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initiazling feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}
