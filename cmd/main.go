package main

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"test/db"
	"test/handlers"
	"test/middleware"
	"test/random"
)

func main() {
	// Загрузка конфигурации
	err := loadConfig()
	if err != nil {
		log.Fatalf("Cannot load config: %s", err)
	}

	// Инициализация логгера
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	// Инициализация приложения Fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(middleware.Authorization)

	// Инициализация базы данных
	db.Init()

	// Роуты
	app.Post("/edit/:Id", handlers.EditNews)
	app.Get("/list", handlers.ListNews)
	app.Post("/login", middleware.Login)

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}

func loadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Генерация случайного секретного ключа, если он не задан в конфиге
	if viper.GetString("JWT_SECRET") == "" {
		secretKey, err := random.GenerateRandomString(32)
		if err != nil {
			return err
		}
		viper.Set("JWT_SECRET", secretKey)
		// Сохраняем обновленную конфигурацию
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}
