package main

import (
	_ "github.com/Mukam21/server_Golang/docs"
	"github.com/Mukam21/server_Golang/pkg/config"
	"github.com/Mukam21/server_Golang/pkg/database"
	"github.com/Mukam21/server_Golang/pkg/handler"
	"github.com/Mukam21/server_Golang/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log := logrus.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	repo := database.NewRepository(db)
	srv := service.NewService(repo, log)
	h := handler.NewHandler(srv, log)

	r := gin.Default()
	h.InitRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info("Starting server on port: ", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
