package controllers

import (
	"net/http"
	"path/filepath"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type articleInput struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	var input articleInput

	file, err := c.FormFile("thumbnail")
	if err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Failed to get file", nil)
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "uploads/articles/"+filename); err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
		return
	}

	articleItem := models.Article{
		ID:           uuid.New(),
		Title:        input.Title,
		Description:  input.Description,
		ThumbnailUrl: "/uploads/articles/" + filename,
		CreatedBy:    userID.(uuid.UUID),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&articleItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create Article", nil)
		return
	}

	utils.RespondSuccess(c, "article created successfully", articleItem)
}

func GetArticles(c *gin.Context) {
	var articles []models.Article

	if err := config.DB.Find(&articles).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to retrieve articles", nil)
		return
	}

	utils.RespondSuccess(c, "articles retrieved successfully", articles)
}

func GetArticleByID(c *gin.Context) {
	articleID := c.Param("id")
	var article models.Article

	if err := config.DB.Where("id = ?", articleID).First(&article).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Article not found", nil)
		return
	}

	utils.RespondSuccess(c, "article retrieved successfully", article)
}
