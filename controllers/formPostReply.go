package controllers

import (
	"net/http"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ForumPostReplyInput struct {
	Description string `form:"description"`
}

func CreateForumPostReply(c *gin.Context) {
	var input ForumPostReplyInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	postID := c.Param("id")
	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "Invalid post ID", nil)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	forumPostReply := models.ForumPostReply{
		ID:          uuid.New(),
		PostID:      parsedPostID,
		Description: input.Description,
		RepliedBy:   userID.(uuid.UUID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&forumPostReply).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create forum post reply", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post reply created successfully", forumPostReply)
}

func GetForumPostReplies(c *gin.Context) {
	postID := c.Param("id")
	var replies []models.ForumPostReply
	if err := config.DB.Where("post_id = ?", postID).Find(&replies).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch forum post replies", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post replies fetched successfully", replies)
}

func GetForumPostReplyByID(c *gin.Context) {
	replyID := c.Param("reply_id")
	var reply models.ForumPostReply
	if err := config.DB.Where("id = ?", replyID).First(&reply).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post reply not found", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post reply fetched successfully", reply)
}

func UpdateForumPostReply(c *gin.Context) {
	var input ForumPostReplyInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	replyID := c.Param("reply_id")
	var forumPostReply models.ForumPostReply

	if err := config.DB.Where("id = ?", replyID).First(&forumPostReply).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post reply not found", nil)
		return
	}

	forumPostReply.Description = input.Description
	forumPostReply.UpdatedAt = time.Now()

	if err := config.DB.Save(&forumPostReply).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update forum post reply", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post reply updated successfully", forumPostReply)
}

func DeleteForumPostReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	var reply models.ForumPostReply
	if err := config.DB.Where("id = ?", replyID).First(&reply).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post reply not found", nil)
		return
	}

	if err := config.DB.Delete(&reply).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to delete forum post reply", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post reply deleted successfully", nil)
}
