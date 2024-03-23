package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mysqlAdapter "gilab.com/estate-agency-api/internal/adapters/db/mysql"
	"gilab.com/estate-agency-api/internal/config"
	"gilab.com/estate-agency-api/internal/domain/service"
	"gilab.com/estate-agency-api/internal/domain/usecase"
	"gilab.com/estate-agency-api/internal/lib/logger"
	"gilab.com/estate-agency-api/internal/storage/mysql"
	v1 "gilab.com/estate-agency-api/internal/transport/http/v1"
	"gilab.com/estate-agency-api/internal/transport/middleware/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.GetLogger("info")

	logger.Info("Connect to db")
	db, err := mysql.New(&cfg.StorageConfig)
	if err != nil {
		panic(err.Error())
	}

	logger.Info("Create router")
	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())
	router.Use(auth.BasicAuth(cfg.HTTPServerConfig.User, cfg.HTTPServerConfig.Password))

	logger.Info("Set routes")
	realtorHandler := v1.NewRealtorHandler(usecase.NewRealtorUsecase(service.NewRealtorService(mysqlAdapter.NewRealtorStorage(db))))
	realtorHandler.Register(router)
	apartmentHandler := v1.NewApartmentHandler(usecase.NewApartmentUsecase(service.NewApartmentService(mysqlAdapter.NewApartmentStorage(db)), service.NewRealtorService(mysqlAdapter.NewRealtorStorage(db))), &logger)
	apartmentHandler.Register(router)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	serv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.Timeout,
		WriteTimeout: cfg.HTTPServerConfig.Timeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	logger.Info("Start server on " + cfg.HTTPServerConfig.Address)

	go func() {
		if err = serv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.TimeClose)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", err.Error())

		return
	}

	if err := db.Close(); err != nil {
		logger.Error("failed to stop db", err.Error())

		return
	}

	logger.Info("server stopped")
}
