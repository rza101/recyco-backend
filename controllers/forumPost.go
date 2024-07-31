package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"recyco/config"
	"recyco/models"
	"recyco/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ForumPostInput struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
}

type ForumPostUpdateInput struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func CreateForumPost(c *gin.Context) {
	var input ForumPostInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, "input tidak valid", nil)
		return
	}

	file, err := c.FormFile("thumbnail")

	var filename string
	if err == nil {
		filename = uuid.New().String() + filepath.Ext(file.Filename)

		if err := c.SaveUploadedFile(file, "uploads/forum/"+filename); err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
			return
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondFailed(c, http.StatusInternalServerError, "User ID not found in context", nil)
		return
	}

	forumPost := models.ForumPost{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		CreatedBy:   userID.(uuid.UUID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if filename != "" {
		forumPost.ThumbnailUrl = "/uploads/forum/" + filename
	}

	if err := config.DB.Create(&forumPost).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to create forum post", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post created successfully", forumPost)
}

func GetForumPosts(c *gin.Context) {
	var posts []models.ForumPost
	if err := config.DB.Preload("PostId").Preload("PostId.RepliedByUser").Preload("CreatedByUser").Find(&posts).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to fetch forum posts", nil)
		return
	}

	var response []map[string]interface{}
	for _, post := range posts {
		postMap := map[string]interface{}{
			"id":            post.ID,
			"title":         post.Title,
			"description":   post.Description,
			"thumbnail_url": post.ThumbnailUrl,
			"created_by":    post.CreatedBy,
			"created_at":    post.CreatedAt,
			"updated_at":    post.UpdatedAt,
			"created_by_user": map[string]interface{}{
				"id":           post.CreatedByUser.ID,
				"phone_number": post.CreatedByUser.PhoneNumber,
				"name":         post.CreatedByUser.Name,
				"role":         post.CreatedByUser.Role,
			},
			"forum_post_replies": []map[string]interface{}{},
		}

		for _, reply := range post.PostId {
			replyMap := map[string]interface{}{
				"id":          reply.ID,
				"post_id":     reply.PostID,
				"description": reply.Description,
				"replied_by":  reply.RepliedBy,
				"created_at":  reply.CreatedAt,
				"updated_at":  reply.UpdatedAt,
				"user": map[string]interface{}{
					"id":           reply.RepliedByUser.ID,
					"phone_number": reply.RepliedByUser.PhoneNumber,
					"name":         reply.RepliedByUser.Name,
					"role":         reply.RepliedByUser.Role,
				},
			}
			postMap["forum_post_replies"] = append(postMap["forum_post_replies"].([]map[string]interface{}), replyMap)
		}
		response = append(response, postMap)
	}

	utils.RespondSuccess(c, "Forum posts fetched successfully", response)
}

func GetForumPostByID(c *gin.Context) {
	var post models.ForumPost
	if err := config.DB.Preload("PostId").Preload("PostId.RepliedByUser").Preload("CreatedByUser").Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post not found", nil)
		return
	}

	postMap := map[string]interface{}{
		"id":            post.ID,
		"title":         post.Title,
		"description":   post.Description,
		"thumbnail_url": post.ThumbnailUrl,
		"created_by":    post.CreatedBy,
		"created_at":    post.CreatedAt,
		"updated_at":    post.UpdatedAt,
		"created_by_user": map[string]interface{}{
			"id":           post.CreatedByUser.ID,
			"phone_number": post.CreatedByUser.PhoneNumber,
			"name":         post.CreatedByUser.Name,
			"role":         post.CreatedByUser.Role,
		},
		"forum_post_replies": []map[string]interface{}{},
	}

	for _, reply := range post.PostId {
		replyMap := map[string]interface{}{
			"id":          reply.ID,
			"post_id":     reply.PostID,
			"description": reply.Description,
			"replied_by":  reply.RepliedBy,
			"created_at":  reply.CreatedAt,
			"updated_at":  reply.UpdatedAt,
			"user": map[string]interface{}{
				"id":           reply.RepliedByUser.ID,
				"phone_number": reply.RepliedByUser.PhoneNumber,
				"name":         reply.RepliedByUser.Name,
				"role":         reply.RepliedByUser.Role,
			},
		}
		postMap["forum_post_replies"] = append(postMap["forum_post_replies"].([]map[string]interface{}), replyMap)
	}

	utils.RespondSuccess(c, "Forum post fetched successfully", postMap)
}

func UpdateForumPost(c *gin.Context) {
	var input ForumPostUpdateInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondFailed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	id := c.Param("id")
	var forumPost models.ForumPost

	if err := config.DB.Where("id = ?", id).First(&forumPost).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post not found", nil)
		return
	}

	file, err := c.FormFile("thumbnail")
	if err == nil {
		if forumPost.ThumbnailUrl != "" {
			os.Remove("." + forumPost.ThumbnailUrl)
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		log.Println("Saving file to:", "uploads/forum/"+filename)

		if err := c.SaveUploadedFile(file, "uploads/forum/"+filename); err != nil {
			utils.RespondFailed(c, http.StatusInternalServerError, "Failed to save file", nil)
			return
		}

		forumPost.ThumbnailUrl = "/uploads/forum/" + filename
	}

	if input.Title != "" {
		forumPost.Title = input.Title
	}
	if input.Description != "" {
		forumPost.Description = input.Description
	}
	forumPost.UpdatedAt = time.Now()

	if err := config.DB.Save(&forumPost).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to update forum post", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post updated successfully", forumPost)
}

func DeleteForumPost(c *gin.Context) {
	var post models.ForumPost
	if err := config.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		utils.RespondFailed(c, http.StatusNotFound, "Forum post not found", nil)
		return
	}

	if post.ThumbnailUrl != "" {
		os.Remove("." + post.ThumbnailUrl)
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		utils.RespondFailed(c, http.StatusInternalServerError, "Failed to delete forum post", nil)
		return
	}

	utils.RespondSuccess(c, "Forum post deleted successfully", nil)
}
