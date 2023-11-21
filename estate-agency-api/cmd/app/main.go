package main

import (
	"net/http"

	"gilab.com/estate-agency-api/internal/adapters/db/mysql"
	"gilab.com/estate-agency-api/internal/config"
	"gilab.com/estate-agency-api/internal/domain/service"
	"gilab.com/estate-agency-api/internal/domain/usecase"
	"gilab.com/estate-agency-api/internal/middlewares"
	v1 "gilab.com/estate-agency-api/internal/transport/http/v1"
	client "gilab.com/estate-agency-api/pkg/client/mysql"
	"gilab.com/estate-agency-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.GetLogger("info")

	logger.Info("Connect to db")
	db, err := client.NewClient(&cfg.StorageConfig)
	if err != nil {
		panic(err)
	}

	logger.Info("Create router")
	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())
	/* routes.UserRoutes(router) */
	router.Use(middlewares.BasicAuth(cfg.HTTPServerConfig.User, cfg.HTTPServerConfig.Password))

	logger.Info("Set routes")
	realtorHandler := v1.NewRealtorHandler(usecase.NewRealtorUsecase(service.NewRealtorService(mysql.NewRealtorStorage(db))))
	realtorHandler.Register(router)
	apartmentHandler := v1.NewApartmentHandler(usecase.NewApartmentUsecase(service.NewApartmentService(mysql.NewApartmentStorage(db)), service.NewRealtorService(mysql.NewRealtorStorage(db))))
	apartmentHandler.Register(router)

	serv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.Timeout,
		WriteTimeout: cfg.HTTPServerConfig.Timeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	logger.Info("Start server on " + cfg.HTTPServerConfig.Address)
	err = serv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
