package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	adapterSql "gilab.com/estate-agency-api/internal/adapters/db/sql"
	"gilab.com/estate-agency-api/internal/config"
	"gilab.com/estate-agency-api/internal/domain/service"
	"gilab.com/estate-agency-api/internal/domain/usecase"
	"gilab.com/estate-agency-api/internal/storage/mysql"
	"gilab.com/estate-agency-api/internal/transport/http/handler"
	"gilab.com/estate-agency-api/internal/transport/http/middleware/auth"
	"github.com/gin-gonic/gin"
)

type app struct {
	cfg *config.Config

	server *http.Server
	db     *sql.DB

	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *app {

	logger.Info("Connect to db")
	db, err := mysql.New(&cfg.StorageConfig)
	if err != nil {
		panic(err)
	}

	logger.Info("Create router")
	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())
	router.Use(auth.BasicAuth(cfg.HTTPServerConfig.User, cfg.HTTPServerConfig.Password))

	logger.Info("Set routes")
	usecase := usecase.NewUsecase(service.NewApartmentService(adapterSql.NewApartmentStorage(db)), service.NewRealtorService(adapterSql.NewRealtorStorage(db)))

	realtorHandler := handler.NewRealtorHandler(usecase, logger)
	realtorHandler.Register(router)

	apartmentHandler := handler.NewApartmentHandler(usecase, logger)
	apartmentHandler.Register(router)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.Timeout,
		WriteTimeout: cfg.HTTPServerConfig.Timeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	return &app{server: server, cfg: cfg, db: db, logger: logger}
}

func (app *app) Run() error {
	return app.server.ListenAndServe()
}

func (app *app) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.cfg.TimeClose)
	defer cancel()

	err := app.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return app.db.Close()
}
