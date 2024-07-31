package controllers

import (
	"net/http"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type TreatmentLocationInput struct {
	Title       string  `form:"title" binding:"required"`
	Address     string  `form:"address" binding:"required"`
	Lat         float64 `form:"lat" binding:"required"`
	Lon         float64 `form:"lon" binding:"required"`
	Description string  `form:"description"`
}

func CreateTreatmentLocation(c *gin.Context) {
	var input TreatmentLocationInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	treatmentLocation := models.TreatmentLocation{
		Title:       input.Title,
		Address:     input.Address,
		Lat:         input.Lat,
		Lon:         input.Lon,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&treatmentLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create treatment location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Treatment location created successfully", "data": treatmentLocation})
}

func GetTreatmentLocations(c *gin.Context) {
	var treatmentLocations []models.TreatmentLocation
	if err := config.DB.Find(&treatmentLocations).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to retrieve treatment locations", nil)
		return
	}

	utils.RespondSuccess(c, "Treatment locations retrieved successfully", treatmentLocations)
}

func GetTreatmentLocationByID(c *gin.Context) {
	treatmentLocationID := c.Param("id")
	var treatmentLocation models.TreatmentLocation

	if err := config.DB.Where("id = ?", treatmentLocationID).First(&treatmentLocation).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Treatment location not found", nil)
		return
	}

	utils.RespondSuccess(c, "Treatment location retrieved successfully", treatmentLocation)
}

func UpdateTreatmentLocation(c *gin.Context) {
	var input TreatmentLocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	treatmentLocationID := c.Param("id")
	var treatmentLocation models.TreatmentLocation

	if err := config.DB.Where("id = ?", treatmentLocationID).First(&treatmentLocation).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Treatment location not found", nil)
		return
	}

	treatmentLocation.Title = input.Title
	treatmentLocation.Address = input.Address
	treatmentLocation.Lat = input.Lat
	treatmentLocation.Lon = input.Lon
	treatmentLocation.Description = input.Description
	treatmentLocation.UpdatedAt = time.Now()

	if err := config.DB.Save(&treatmentLocation).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update treatment location", nil)
		return
	}

	utils.RespondSuccess(c, "Treatment location updated successfully", treatmentLocation)
}

func DeleteTreatmentLocation(c *gin.Context) {
	treatmentLocationID := c.Param("id")
	var treatmentLocation models.TreatmentLocation

	if err := config.DB.Where("id = ?", treatmentLocationID).First(&treatmentLocation).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Treatment location not found", nil)
		return
	}

	if err := config.DB.Delete(&treatmentLocation).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to delete treatment location", nil)
		return
	}

	utils.RespondSuccess(c, "Treatment location deleted successfully", nil)
}
