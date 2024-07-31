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

type communityInput struct {
	Name         string `form:"name" binding:"required"`
	Description  string `form:"description"`
	CommunityUrl string `form:"community_url" binding:"required"`
}

func CreateCommunity(c *gin.Context) {
	var input communityInput

	file, err := c.FormFile("thumbnail")
	if err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Failed to get file", nil)
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "uploads/community/"+filename); err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
		return
	}

	communityItem := models.Community{
		ID:           uuid.New(),
		Name:         input.Name,
		Description:  input.Description,
		Community:    input.CommunityUrl,
		ThumbnailUrl: "/uploads/community/" + filename,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&communityItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create Community", nil)
		return
	}

	utils.RespondSuccess(c, "community created successfully", communityItem)
}

func GetCommunities(c *gin.Context) {
	var communities []models.Community
	if err := config.DB.Find(&communities).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to retrieve communities", nil)
		return
	}

	utils.RespondSuccess(c, "Communities retrieved successfully", communities)
}
