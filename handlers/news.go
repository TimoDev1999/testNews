package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/sirupsen/logrus"
	"strconv"
	"test/db"
	"test/models"
)

type EditNewsRequest struct {
	Title      *string `json:"Title" validate:"omitempty,max=255"`
	Content    *string `json:"Content" validate:"omitempty"`
	Categories []int64 `json:"Categories" validate:"omitempty,dive,gt=0"`
}

var validate = validator.New()

func EditNews(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("Id"), 10, 64)
	if err != nil {
		logrus.Error("Invalid news ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid news ID"})
	}

	var request EditNewsRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.Error("Cannot parse JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validate.Struct(&request); err != nil {
		logrus.Error("Validation error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tx, err := db.DB.Begin()
	if err != nil {
		logrus.Error("Database transaction error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database transaction error"})
	}
	defer tx.Rollback()

	news := &models.News{ID: id}
	err = tx.Model(news).WherePK().Select()
	if err != nil {
		logrus.Error("News not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
	}

	if request.Title != nil {
		news.Title = *request.Title
	}
	if request.Content != nil {
		news.Content = *request.Content
	}

	_, err = tx.Model(news).WherePK().UpdateNotZero()
	if err != nil {
		logrus.Error("Failed to update news:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
	}

	if request.Categories != nil {
		_, err = tx.Model(&models.NewsCategories{}).Where("news_id = ?", id).Delete()
		if err != nil {
			logrus.Error("Failed to update categories:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update categories"})
		}

		for _, categoryId := range request.Categories {
			newsCategory := &models.NewsCategories{NewsID: id, CategoryID: categoryId}
			_, err := tx.Model(newsCategory).Insert()
			if err != nil {
				logrus.Error("Failed to update categories:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update categories"})
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error("Database commit error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database commit error"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func ListNews(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var newsList []models.News
	err := db.DB.Model(&newsList).Limit(limit).Offset(offset).Select()
	if err != nil {
		logrus.Error("Failed to fetch news:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}

	response := make([]map[string]interface{}, len(newsList))
	for i, news := range newsList {
		var categories []int64
		err := db.DB.Model((*models.NewsCategories)(nil)).Column("category_id").Where("news_id = ?", news.ID).Select(&categories)
		if err != nil {
			logrus.Error("Failed to fetch categories:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch categories"})
		}
		response[i] = map[string]interface{}{
			"Id":         news.ID,
			"Title":      news.Title,
			"Content":    news.Content,
			"Categories": categories,
		}
	}

	return c.JSON(fiber.Map{"Success": true, "News": response})
}
