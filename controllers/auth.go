package controllers

import (
	"net/http"
	"recyco/config"
	"recyco/middlewares"
	"recyco/models"
	"recyco/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	PhoneNumber string `form:"phone_number" binding:"required"`
	Password    string `form:"password" binding:"required"`
	Name        string `form:"name" binding:"required"`
	Role        string `form:"role" binding:"required"`
}

type LoginInput struct {
	PhoneNumber string `form:"phone_number" binding:"required"`
	Password    string `form:"password" binding:"required"`
}

type UpdateInput struct {
	PhoneNumber string `form:"phone_number" binding:"required"`
	Password    string `form:"password"`
	Name        string `form:"name" binding:"required"`
	Role        string `form:"role" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	if len(input.Password) < 8 || len(input.Password) > 50 {
		utils.RespondFailed(c, http.StatusBadRequest, "Password must be between 8 and 50 characters", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Password encryption failed", nil)
		return
	}

	user := models.User{
		PhoneNumber: input.PhoneNumber,
		Password:    string(hashedPassword),
		Name:        input.Name,
		Role:        input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {

		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062") {
			utils.RespondFailed(c, http.StatusConflict, "Nomor telepon ini sudah terdaftar. Silakan gunakan nomor lain.", nil)
			return
		}
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	responseData := gin.H{
		"id":           user.ID,
		"phone_number": user.PhoneNumber,
		"name":         user.Name,
		"role":         user.Role,
	}

	utils.RespondSuccess(c, "User registered successfully", responseData)
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var user models.User
	if err := config.DB.Where("phone_number = ?", input.PhoneNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondFailed(c, http.StatusUnauthorized, "Phone number not found", nil)
			return
		}
		utils.RespondFailed(c, http.StatusInternalServerError, "Database error", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.RespondFailed(c, http.StatusUnauthorized, "Incorrect password", nil)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	utils.RespondSuccess(c, "Successfully logged in", gin.H{"token": token})
}

func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "User not found", nil)
		return
	}

	responseData := gin.H{
		"id":           user.ID,
		"phone_number": user.PhoneNumber,
		"name":         user.Name,
		"role":         user.Role,
	}

	utils.RespondSuccess(c, "User profile retrieved successfully", responseData)
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.RespondFailed(c, http.StatusUnauthorized, "Authorization header required", nil)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.RespondFailed(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
		return
	}

	token := parts[1]
	middlewares.AddToBlacklist(token)

	utils.RespondSuccess(c, "Successfully logged out", nil)
}
