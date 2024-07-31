package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketItemInput struct {
	Name        string  `form:"name" binding:"required"`
	Price       float64 `form:"price" binding:"required"`
	Weight      float64 `form:"weight" binding:"required"`
	Description string  `form:"description"`
	OrderedBy   string  `form:"ordered_by"`
}

type UpdateMarketItemInput struct {
	Name        string  `form:"name"`
	Price       float64 `form:"price"`
	Weight      float64 `form:"weight"`
	Description string  `form:"description"`
	OrderedBy   string  `form:"ordered_by"`
	Status      string  `form:"status"`
}

func CreateMarketItem(c *gin.Context) {
	var input MarketItemInput

	file, err := c.FormFile("thumbnail")
	if err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Failed to get file", nil)
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

	userRole, exists := c.Get("userRole")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User role not found in context", nil)
		return
	}

	var itemScale string
	if userRole == "P_SMALL" {
		itemScale = "SMALL"
	} else if userRole == "P_LARGE" {
		itemScale = "LARGE"
	} else {
		utils.RespondFailed(c, http.StatusForbidden, "Invalid user role for creating market item", nil)
		return
	}

	if itemScale == "SMALL" && input.Weight > 15 {
		utils.RespondFailed(c, http.StatusNotFound, "Weight for SMALL scale items must not exceed 15", nil)
		return
	}

	if itemScale == "LARGE" && input.Weight <= 15 {
		utils.RespondFailed(c, http.StatusNotFound, "Weight for LARGE scale items must be greater than 15", nil)
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "uploads/markets/"+filename); err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
		return
	}

	marketItem := models.MarketItems{
		ID:           uuid.New(),
		Name:         input.Name,
		Price:        input.Price,
		Weight:       input.Weight,
		ItemScale:    itemScale,
		Description:  input.Description,
		ThumbnailUrl: "/uploads/markets/" + filename,
		PostedBy:     userID.(uuid.UUID),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create market item", nil)
		return
	}

	response := map[string]interface{}{
		"id":             marketItem.ID,
		"name":           marketItem.Name,
		"price":          marketItem.Price,
		"weight":         marketItem.Weight,
		"item_scale":     marketItem.ItemScale,
		"description":    marketItem.Description,
		"thumbnail_url":  marketItem.ThumbnailUrl,
		"posted_by":      marketItem.PostedBy,
		"ordered_by":     input.OrderedBy,
		"created_at":     marketItem.CreatedAt,
		"updated_at":     marketItem.UpdatedAt,
		"deleted_at":     marketItem.DeletedAt,
		"posted_by_user": userID,
	}

	utils.RespondSuccess(c, "Market item created successfully", response)
}

func GetMarketItems(c *gin.Context) {
	var items []models.MarketItems

	userRole, exists := c.Get("userRole")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User role not found in context", nil)
		return
	}
	switch userRole {
	case "C_SMALL":
		if err := config.DB.Preload("PostedByUser").Where("item_scale = ?", "SMALL").Find(&items).Error; err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
			return
		}
	case "C_LARGE":
		if err := config.DB.Preload("PostedByUser").Where("item_scale = ?", "LARGE").Find(&items).Error; err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
			return
		}
	case "P_SMALL":
		if err := config.DB.Preload("PostedByUser").Where("item_scale = ?", "SMALL").Find(&items).Error; err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
			return
		}
	case "P_LARGE":
		if err := config.DB.Preload("PostedByUser").Where("item_scale = ?", "LARGE").Find(&items).Error; err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
			return
		}
	default:
		if err := config.DB.Preload("PostedByUser").Find(&items).Error; err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
			return
		}
	}

	if len(items) == 0 {
		utils.RespondSuccess(c, "Market items fetched successfully", []map[string]interface{}{})
		return
	}

	var responseItems []map[string]interface{}
	for _, item := range items {
		postedBy := map[string]interface{}{
			"id":           item.PostedByUser.ID,
			"name":         item.PostedByUser.Name,
			"phone_number": item.PostedByUser.PhoneNumber,
			"role":         item.PostedByUser.Role,
		}

		responseItem := map[string]interface{}{
			"id":            item.ID,
			"name":          item.Name,
			"price":         item.Price,
			"weight":        item.Weight,
			"scale":         item.ItemScale,
			"description":   item.Description,
			"thumbnail_url": item.ThumbnailUrl,
			"posted_by":     postedBy,
			"created_at":    item.CreatedAt,
			"updated_at":    item.UpdatedAt,
			"deleted_at":    item.DeletedAt,
		}
		responseItems = append(responseItems, responseItem)
	}

	utils.RespondSuccess(c, "Market items fetched successfully", responseItems)
}

func GetMarketItemByID(c *gin.Context) {
	var item models.MarketItems

	itemID := c.Param("id")

	if err := config.DB.Preload("PostedByUser").Where("id = ?", itemID).First(&item).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User role not found in context", nil)
		return
	}

	if (userRole == "C_SMALL" && item.ItemScale != "SMALL") || (userRole == "C_LARGE" && item.ItemScale != "LARGE") {
		utils.RespondFailed(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
		return
	}

	postedBy := map[string]interface{}{
		"id":           item.PostedByUser.ID,
		"name":         item.PostedByUser.Name,
		"phone_number": item.PostedByUser.PhoneNumber,
		"role":         item.PostedByUser.Role,
	}

	responseItem := map[string]interface{}{
		"id":            item.ID,
		"name":          item.Name,
		"price":         item.Price,
		"weight":        item.Weight,
		"description":   item.Description,
		"thumbnail_url": item.ThumbnailUrl,
		"posted_by":     postedBy,
		"created_at":    item.CreatedAt,
		"updated_at":    item.UpdatedAt,
		"deleted_at":    item.DeletedAt,
	}

	utils.RespondSuccess(c, "Market item fetched successfully", responseItem)
}

func GetUserMarketItems(c *gin.Context) {
	var items []models.MarketItems

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	if err := config.DB.Preload("PostedByUser").Unscoped().Where("posted_by = ?", userID).Find(&items).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market items", nil)
		return
	}

	if len(items) == 0 {
		utils.RespondSuccess(c, "User market items fetched successfully", []map[string]interface{}{})
		return
	}

	var responseItems []map[string]interface{}
	for _, item := range items {
		postedBy := map[string]interface{}{
			"id":           item.PostedByUser.ID,
			"name":         item.PostedByUser.Name,
			"phone_number": item.PostedByUser.PhoneNumber,
			"role":         item.PostedByUser.Role,
		}

		var transaction models.MarketItemTransactionActivities
		status := "READY"

		if err := config.DB.Unscoped().Where("item_id = ?", item.ID).Order("created_at desc").First(&transaction).Error; err == nil {
			status = transaction.Status
		}

		responseItem := map[string]interface{}{
			"id":            item.ID,
			"name":          item.Name,
			"price":         item.Price,
			"weight":        item.Weight,
			"scale":         item.ItemScale,
			"description":   item.Description,
			"thumbnail_url": item.ThumbnailUrl,
			"posted_by":     postedBy,
			"created_at":    item.CreatedAt,
			"updated_at":    item.UpdatedAt,
			"deleted_at":    item.DeletedAt,
			"status":        status,
		}
		responseItems = append(responseItems, responseItem)
	}

	utils.RespondSuccess(c, "User market items fetched successfully", responseItems)
}

func UpdateMarketItem(c *gin.Context) {
	var input UpdateMarketItemInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Input tidak valid", nil)
		return
	}

	id := c.Param("id")
	var marketItem models.MarketItems

	if err := config.DB.Where("id = ?", id).First(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User role not found in context", nil)
		return
	}

	if marketItem.PostedBy != userID {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to update this market item", nil)
		return
	}

	if input.Weight > 0 {
		if userRole == "P_SMALL" && input.Weight > 15 {
			utils.RespondFailed(c, http.StatusForbidden, "Weight for SMALL scale items must not exceed 15", nil)
			return
		}

		if userRole == "P_LARGE" && input.Weight <= 15 {
			utils.RespondFailed(c, http.StatusForbidden, "Weight for LARGE scale items must be greater than 15", nil)
			return
		}
	}

	if userRole == "P_SMALL" && marketItem.ItemScale != "SMALL" {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to update LARGE scale items", nil)
		return
	}
	if userRole == "P_LARGE" && marketItem.ItemScale != "LARGE" {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to update SMALL scale items", nil)
		return
	}

	file, err := c.FormFile("thumbnail")
	if err == nil {
		if marketItem.ThumbnailUrl != "" {
			os.Remove("." + marketItem.ThumbnailUrl)
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, "uploads/market/"+filename); err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
			return
		}

		marketItem.ThumbnailUrl = "/uploads/market/" + filename
	}

	if input.Name != "" {
		marketItem.Name = input.Name
	}
	if input.Price != 0 {
		marketItem.Price = input.Price
	}
	if input.Weight != 0 {
		marketItem.Weight = input.Weight
	}
	if input.Description != "" {
		marketItem.Description = input.Description
	}

	marketItem.UpdatedAt = time.Now()

	tx := config.DB.Begin()
	if tx.Error != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	if input.Status == "FINISHED" {
		transactionActivity := models.MarketItemTransactionActivities{
			ID:              uuid.New(),
			ItemID:          marketItem.ID,
			Status:          input.Status,
			TransactionByID: uuid.Nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		var pickupInfo models.MarketItemPickupInformations
		if err := config.DB.Where("item_id = ?", marketItem.ID).First(&pickupInfo).Error; err == nil {
			transactionActivity.TransactionByID = pickupInfo.ID
		}

		if err := tx.Create(&transactionActivity).Error; err != nil {
			tx.Rollback()
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create transaction activity", nil)
			return
		}

		marketItem.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	}

	if err := tx.Save(&marketItem).Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update market item", nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	var transaction models.MarketItemTransactionActivities
	status := "READY"

	if err := config.DB.Unscoped().Where("item_id = ?", marketItem.ID).Order("created_at desc").First(&transaction).Error; err == nil {
		status = transaction.Status
	}

	response := map[string]interface{}{
		"id":            marketItem.ID,
		"name":          marketItem.Name,
		"price":         marketItem.Price,
		"weight":        marketItem.Weight,
		"item_scale":    marketItem.ItemScale,
		"description":   marketItem.Description,
		"thumbnail_url": marketItem.ThumbnailUrl,
		"posted_by":     marketItem.PostedBy,
		"ordered_by":    input.OrderedBy,
		"created_at":    marketItem.CreatedAt,
		"updated_at":    marketItem.UpdatedAt,
		"deleted_at":    marketItem.DeletedAt,
		"status":        status,
	}

	utils.RespondSuccess(c, "Market item updated successfully", response)
}

func DeleteMarketItem(c *gin.Context) {
	id := c.Param("id")
	var marketItem models.MarketItems

	if err := config.DB.Where("id = ?", id).First(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User role not found in context", nil)
		return
	}

	if marketItem.PostedBy != userID {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to delete this market item", nil)
		return
	}

	if userRole == "P_SMALL" && marketItem.ItemScale != "SMALL" {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to delete LARGE scale items", nil)
		return
	}
	if userRole == "P_LARGE" && marketItem.ItemScale != "LARGE" {
		utils.RespondFailed(c, http.StatusForbidden, "You are not allowed to delete SMALL scale items", nil)
		return
	}

	if marketItem.ThumbnailUrl != "" {
		os.Remove("." + marketItem.ThumbnailUrl)
	}

	if err := config.DB.Delete(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to delete market item", nil)
		return
	}

	utils.RespondSuccess(c, "Market item deleted successfully", nil)
}
