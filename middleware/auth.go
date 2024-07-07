package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET"))

// Authorization проверяет наличие и валидность JWT токена в заголовке Authorization
func Authorization(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	// Проверяем наличие токена
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: token is missing",
		})
	}

	// Парсим и проверяем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Unauthorized: %v", err),
		})
	}

	// Проверяем, что токен верный
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: token is invalid",
		})
	}

	return c.Next()
}

// Login обрабатывает запрос на аутентификацию и генерирует JWT токен при успешной аутентификации
func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Username != "user" || req.Password != "password" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = req.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Токен действителен 24 часа

	// Подписываем токен с помощью секретного ключа
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT token",
		})
	}

	// Возвращаем токен как ответ
	return c.JSON(fiber.Map{
		"token": signedToken,
	})
}
