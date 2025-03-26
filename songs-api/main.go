package main

import (
	"fmt"
	"songs-api/config"
	"songs-api/handlers"
	"songs-api/models" // Добавляем импорт для SongDetail
	"songs-api/storage"

	_ "songs-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	cfg := config.LoadConfig()

	// Подключение к базе данных и миграции
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации миграций: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Ошибка при выполнении миграций: %v", err)
	}
	logrus.Info("Миграции успешно выполнены")

	// Инициализация хранилища
	store, err := storage.NewStorage(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Настройка маршрутов
	r := gin.Default()
	r.GET("/songs", handlers.GetSongs(store))
	r.POST("/songs", handlers.AddSong(cfg, store))
	r.GET("/songs/:id", handlers.GetSongText(store))
	r.PUT("/songs/:id", handlers.UpdateSong(store))
	r.DELETE("/songs/:id", handlers.DeleteSong(store))

	// Новый маршрут /info
	r.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")
		c.JSON(200, models.SongDetail{
			ReleaseDate: "2023-01-01", // Заглушка, можно заменить на реальные данные
			Text:        "Текст песни для " + group + " - " + song,
			Link:        "http://example.com/" + group + "/" + song,
		})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
