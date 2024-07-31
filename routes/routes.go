package routes

import (
	"recyco/controllers"
	"recyco/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	imageRoutes := r.Group("/uploads")
	{
		imageRoutes.Static("/forum", "uploads/forum")
		imageRoutes.Static("/markets", "uploads/markets")
		imageRoutes.Static("/articles", "uploads/articles")
		imageRoutes.Static("/community", "uploads/community")
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/logout", controllers.Logout)
	}

	articles := r.Group("/articles")
	articles.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		articles.POST("/", middlewares.RoleMiddleware("ADMIN"), controllers.CreateArticle)
		articles.GET("/", controllers.GetArticles)
	}

	treatmentLocations := r.Group("/treatment_locations")
	treatmentLocations.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		treatmentLocations.POST("/", controllers.CreateTreatmentLocation)
		treatmentLocations.GET("/", controllers.GetTreatmentLocations)
		treatmentLocations.GET("/:id", controllers.GetTreatmentLocationByID)
		treatmentLocations.PUT("/:id", controllers.UpdateTreatmentLocation)
		treatmentLocations.DELETE("/:id", controllers.DeleteTreatmentLocation)
	}

	userRoutes := r.Group("/user")
	userRoutes.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		userRoutes.GET("/profile", controllers.GetUserProfile)
	}

	marketItems := r.Group("/markets")
	marketItems.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		marketItems.POST("/", middlewares.RoleMiddleware("ADMIN", "P_SMALL", "P_LARGE"), controllers.CreateMarketItem)
		marketItems.GET("/", controllers.GetMarketItems)
		marketItems.GET("/:id", controllers.GetMarketItemByID)
		marketItems.GET("/markets_self", controllers.GetUserMarketItems)
		marketItems.PUT("/:id", middlewares.RoleMiddleware("P_SMALL", "P_LARGE"), controllers.UpdateMarketItem)
		marketItems.DELETE("/:id", middlewares.RoleMiddleware("P_SMALL", "P_LARGE"), controllers.DeleteMarketItem)
	}
	marketItemsSelf := r.Group("/markets_self")
	marketItemsSelf.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		marketItemsSelf.GET("/", controllers.GetUserMarketItems)
	}

	marketTransactions := r.Group("/market_transactions")
	marketTransactions.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		marketTransactions.GET("/", middlewares.RoleMiddleware("P_LARGE", "C_LARGE"), controllers.GetMarketItemTransactions)
		marketTransactions.POST("/", middlewares.RoleMiddleware("C_LARGE"), controllers.CreateMarketItemPickupInformation)
		marketTransactions.GET("/:id", middlewares.RoleMiddleware("P_LARGE", "C_LARGE"), controllers.GetMarketItemPickupInformationByID)
		marketTransactions.PUT("/:id", middlewares.RoleMiddleware("P_LARGE"), controllers.UpdateMarketItemTransactionStatus)
	}

	communityRoutes := r.Group("/communities")
	communityRoutes.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		communityRoutes.POST("/", middlewares.RoleMiddleware("ADMIN"), controllers.CreateCommunity)
		communityRoutes.GET("/", controllers.GetCommunities)
	}

	forumPosts := r.Group("/forum_posts")
	forumPosts.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		forumPosts.POST("/", controllers.CreateForumPost)
		forumPosts.GET("/", controllers.GetForumPosts)
		forumPosts.GET("/:id", controllers.GetForumPostByID)
		forumPosts.PUT("/:id", controllers.UpdateForumPost)
		forumPosts.DELETE("/:id", controllers.DeleteForumPost)
	}

	forumPostReplies := r.Group("/forum_posts/:id/replies")
	forumPostReplies.Use(middlewares.JWTAuthMiddleware(), middlewares.JWTBlacklistMiddleware())
	{
		forumPostReplies.GET("/", controllers.GetForumPostReplies)
		forumPostReplies.POST("/", controllers.CreateForumPostReply)
		forumPostReplies.PUT("/:reply_id", controllers.UpdateForumPostReply)
		forumPostReplies.DELETE("/:reply_id", controllers.DeleteForumPostReply)
	}

	return r
}
