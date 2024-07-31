package controllers

import (
	"net/http"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketItemPickupInput struct {
	ItemID                    string `form:"item_id" binding:"required"`
	RecipientName             string `form:"recipient_name" binding:"required"`
	RecipientPhone            string `form:"recipient_phone" binding:"required"`
	Description               string `form:"description"`
	PickupLocationAddress     string `form:"pickup_location_address" binding:"required"`
	PickupLocationDescription string `form:"pickup_location_description"`
	Status                    string `form:"status"`
}

type MarketItemTransactionStatusUpdateInput struct {
	Status string `form:"status" binding:"required"`
}

func CreateMarketItemPickupInformation(c *gin.Context) {
	var input MarketItemPickupInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Input tidak valid", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists || userRole != "C_LARGE" {
		utils.RespondFailed(c, http.StatusForbidden, "You do not have permission to create pickup information", nil)
		return
	}

	if input.Status == "" {
		input.Status = "ON_PROCESS"
	}

	validStatuses := []string{"ON_PROCESS", "ON_DELIVER", "FINISHED", "CANCELLED"}
	isValidStatus := false
	for _, status := range validStatuses {
		if input.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		utils.RespondFailed(c, http.StatusBadRequest, "Invalid status value", nil)
		return
	}

	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Invalid UUID for item_id", nil)
		return
	}

	var marketItem models.MarketItems
	if err := config.DB.Where("id = ?", itemID).First(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	var servicePrice int
	if marketItem.Weight <= 50 {
		servicePrice = 1000
	} else {
		servicePrice = (int(marketItem.Weight/50) + 1) * 1000
	}
	deliveryPrice := marketItem.Weight * 1500

	tx := config.DB.Begin()
	if tx.Error != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	pickupInformation := models.MarketItemPickupInformations{
		ID:                        uuid.New(),
		ItemID:                    itemID,
		RecipientName:             input.RecipientName,
		RecipientPhone:            input.RecipientPhone,
		Description:               input.Description,
		PickupLocationAddress:     input.PickupLocationAddress,
		PickupLocationDescription: input.PickupLocationDescription,
		ServicePrice:              float64(servicePrice),
		DeliveryPrice:             deliveryPrice,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	if err := tx.Create(&pickupInformation).Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create pickup information", nil)
		return
	}

	transactionActivity := models.MarketItemTransactionActivities{
		ID:              uuid.New(),
		ItemID:          itemID,
		Status:          input.Status,
		TransactionByID: pickupInformation.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&transactionActivity).Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create transaction activity", nil)
		return
	}

	if input.Status == "ON_PROCESS" {
		if err := tx.Model(&marketItem).Updates(map[string]interface{}{
			"DeletedAt": time.Now(),
			"OrderedBy": userID,
		}).Error; err != nil {
			tx.Rollback()
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update market item", nil)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	utils.RespondSuccess(c, "Pickup information and transaction activity created successfully", gin.H{
		"pickup_information": gin.H{
			"id":                          pickupInformation.ID,
			"item_id":                     pickupInformation.ItemID,
			"recipient_name":              pickupInformation.RecipientName,
			"recipient_phone":             pickupInformation.RecipientPhone,
			"description":                 pickupInformation.Description,
			"pickup_location_address":     pickupInformation.PickupLocationAddress,
			"pickup_location_description": pickupInformation.PickupLocationDescription,
			"service_price":               pickupInformation.ServicePrice,
			"delivery_price":              pickupInformation.DeliveryPrice,
			"created_at":                  pickupInformation.CreatedAt,
			"updated_at":                  pickupInformation.UpdatedAt,
		},
		"transaction_activity": gin.H{
			"id":             transactionActivity.ID,
			"item_id":        transactionActivity.ItemID,
			"status":         transactionActivity.Status,
			"created_at":     transactionActivity.CreatedAt,
			"updated_at":     transactionActivity.UpdatedAt,
			"transaction_by": transactionActivity.TransactionByID,
		},
	})
}

func GetMarketItemTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var latestTransactions []models.MarketItemTransactionActivities
	subQuery := config.DB.Unscoped().Table("market_item_transaction_activities").
		Select("MAX(created_at)").
		Group("item_id")

	if err := config.DB.Unscoped().Preload("MarketItem.PostedByUser").
		Preload("MarketItem.OrderedByUser").
		Where("created_at IN (?)", subQuery).
		Find(&latestTransactions).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch market transactions", nil)
		return
	}

	var responseItems []map[string]interface{}
	for _, transaction := range latestTransactions {
		marketItem := transaction.MarketItem
		if marketItem.PostedBy != userID && (marketItem.OrderedBy == nil || *marketItem.OrderedBy != userID) {
			continue
		}

		postedBy := map[string]interface{}{
			"id":           marketItem.PostedByUser.ID,
			"phone_number": marketItem.PostedByUser.PhoneNumber,
			"name":         marketItem.PostedByUser.Name,
			"role":         marketItem.PostedByUser.Role,
		}

		var orderedBy map[string]interface{}
		if marketItem.OrderedBy != nil {
			orderedBy = map[string]interface{}{
				"id":           marketItem.OrderedByUser.ID,
				"phone_number": marketItem.OrderedByUser.PhoneNumber,
				"name":         marketItem.OrderedByUser.Name,
				"role":         marketItem.OrderedByUser.Role,
			}
		} else {
			orderedBy = nil
		}

		item := map[string]interface{}{
			"id":            marketItem.ID,
			"name":          marketItem.Name,
			"price":         marketItem.Price,
			"weight":        marketItem.Weight,
			"description":   marketItem.Description,
			"thumbnail_url": marketItem.ThumbnailUrl,
			"posted_by":     postedBy,
			"posted_at":     marketItem.CreatedAt.Format(time.RFC3339),
			"ordered_by":    orderedBy,
		}

		responseItem := map[string]interface{}{
			"id":   transaction.ID,
			"item": item,
			"last_status": map[string]interface{}{
				"status":     transaction.Status,
				"created_at": transaction.CreatedAt.Format(time.RFC3339),
			},
			"transaction_by": transaction.TransactionByID,
		}

		responseItems = append(responseItems, responseItem)
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Market transactions fetched successfully",
		"data":    responseItems,
	}

	if len(responseItems) == 0 {
		response["data"] = []interface{}{}
	}

	c.JSON(http.StatusOK, response)
}

func GetMarketItemPickupInformationByID(c *gin.Context) {
	itemID := c.Param("id")
	var marketItem models.MarketItems

	if err := config.DB.Unscoped().Preload("PostedByUser").Preload("OrderedByUser").Where("id = ?", itemID).First(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	var pickupInfo models.MarketItemPickupInformations
	if err := config.DB.Preload("MarketItem").Where("item_id = ?", marketItem.ID).First(&pickupInfo).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Pickup information not found", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if marketItem.PostedBy != userID && (marketItem.OrderedBy == nil || *marketItem.OrderedBy != userID) {
		utils.RespondFailed(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
		return
	}

	var transactionActivities []models.MarketItemTransactionActivities
	if err := config.DB.Unscoped().Where("item_id = ?", marketItem.ID).Order("created_at desc").Find(&transactionActivities).Error; err != nil {
		transactionActivities = []models.MarketItemTransactionActivities{}
	}

	postedBy := map[string]interface{}{
		"id":           marketItem.PostedByUser.ID,
		"phone_number": marketItem.PostedByUser.PhoneNumber,
		"name":         marketItem.PostedByUser.Name,
		"role":         marketItem.PostedByUser.Role,
	}

	orderedBy := map[string]interface{}{}
	if marketItem.OrderedBy != nil {
		orderedBy = map[string]interface{}{
			"id":           marketItem.OrderedByUser.ID,
			"phone_number": marketItem.OrderedByUser.PhoneNumber,
			"name":         marketItem.OrderedByUser.Name,
			"role":         marketItem.OrderedByUser.Role,
		}
	}

	itemDetail := map[string]interface{}{
		"id":            marketItem.ID,
		"name":          marketItem.Name,
		"price":         marketItem.Price,
		"weight":        marketItem.Weight,
		"description":   marketItem.Description,
		"thumbnail_url": marketItem.ThumbnailUrl,
		"posted_by":     postedBy,
		"ordered_by":    orderedBy,
		"posted_at":     marketItem.CreatedAt.Format(time.RFC3339),
	}

	var allStatus []map[string]interface{}
	for _, activity := range transactionActivities {
		if marketItem.OrderedBy != nil && activity.TransactionByID == pickupInfo.ID {
			allStatus = append(allStatus, map[string]interface{}{
				"status":     activity.Status,
				"created_at": activity.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	if allStatus == nil {
		allStatus = []map[string]interface{}{}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Pickup information fetched successfully",
		"data": map[string]interface{}{
			"id":         pickupInfo.ID,
			"item":       itemDetail,
			"all_status": allStatus,
		},
	}

	c.JSON(http.StatusOK, response)
}

func UpdateMarketItemTransactionStatus(c *gin.Context) {
	var input MarketItemTransactionStatusUpdateInput
	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Input tidak valid", nil)
		return
	}

	itemID := c.Param("id")
	var marketItem models.MarketItems

	if err := config.DB.Unscoped().Where("id = ?", itemID).First(&marketItem).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Market item not found", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if marketItem.PostedBy != userID {
		utils.RespondFailed(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
		return
	}

	var pickupInfo models.MarketItemPickupInformations
	if err := config.DB.Where("item_id = ?", marketItem.ID).First(&pickupInfo).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Pickup information not found", nil)
		return
	}

	validStatuses := []string{"ON_PROCESS", "ON_DELIVER", "FINISHED", "CANCELLED"}
	isValidStatus := false
	for _, status := range validStatuses {
		if input.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		utils.RespondFailed(c, http.StatusBadRequest, "Invalid status value", nil)
		return
	}

	var existingTransaction models.MarketItemTransactionActivities
	if err := config.DB.Where("item_id = ? AND status = ? AND transaction_by_id = ?", marketItem.ID, input.Status, pickupInfo.ID).First(&existingTransaction).Error; err == nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Duplicate status for the same item and transaction", nil)
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	transaction := models.MarketItemTransactionActivities{
		ID:              uuid.New(),
		ItemID:          marketItem.ID,
		Status:          input.Status,
		TransactionByID: pickupInfo.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create new transaction activity", nil)
		return
	}

	if input.Status == "CANCELLED" {
		marketItem.OrderedBy = nil
		marketItem.DeletedAt = gorm.DeletedAt{}
		if err := tx.Save(&marketItem).Error; err != nil {
			tx.Rollback()
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update market item", nil)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	response := map[string]interface{}{
		"id":             transaction.ID,
		"item_id":        transaction.ItemID,
		"status":         transaction.Status,
		"created_at":     transaction.CreatedAt,
		"updated_at":     transaction.UpdatedAt,
		"transaction_by": transaction.TransactionByID,
	}

	utils.RespondSuccess(c, "Transaction status updated successfully", response)
}
