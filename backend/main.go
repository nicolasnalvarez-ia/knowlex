package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twitter-bookmarks-api/auth"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/handlers"
	"twitter-bookmarks-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	auth.InitOAuth()

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	
	// Allow frontend URL and Chrome extension origins
	corsConfig.AllowOriginFunc = func(origin string) bool {
		// Allow frontend URL
		if origin == frontendURL {
			return true
		}
		// Allow Chrome extension origins
		if len(origin) > 16 && origin[:16] == "chrome-extension" {
			return true
		}
		return false
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.Logger())

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.GET("/twitter", handlers.TwitterAuth)
			authGroup.GET("/twitter/callback", handlers.TwitterCallback)
			authGroup.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
			authGroup.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
		}

		bookmarksGroup := api.Group("/bookmarks")
		bookmarksGroup.Use(middleware.AuthMiddleware())
		{
			bookmarksGroup.GET("", handlers.GetBookmarks)
			bookmarksGroup.POST("/import", handlers.ImportBookmarks)
			bookmarksGroup.DELETE("/:id", handlers.DeleteBookmark)
			bookmarksGroup.GET("/search", handlers.SearchBookmarks)
			bookmarksGroup.POST("/:id/category", handlers.AssignCategory)
			bookmarksGroup.DELETE("/:id/category/:categoryId", handlers.RemoveCategory)
		}

		categoriesGroup := api.Group("/categories")
		categoriesGroup.Use(middleware.AuthMiddleware())
		{
			categoriesGroup.GET("", handlers.GetCategories)
			categoriesGroup.POST("", handlers.CreateCategory)
			categoriesGroup.PUT("/:id", handlers.UpdateCategory)
			categoriesGroup.DELETE("/:id", handlers.DeleteCategory)
		}

		exportGroup := api.Group("/export")
		exportGroup.Use(middleware.AuthMiddleware())
		{
			exportGroup.GET("/bookmarks", handlers.ExportBookmarks)
			exportGroup.GET("/category/:id", handlers.ExportCategory)
		}

		userGroup := api.Group("/user")
		userGroup.Use(middleware.AuthMiddleware())
		{
			userGroup.DELETE("/account", handlers.DeleteAccount)
		}

		aiGroup := api.Group("/ai")
		aiGroup.Use(middleware.AuthMiddleware())
		{
			aiGroup.POST("/categorize", handlers.AutoCategorizeBookmarks)
			aiGroup.POST("/categorize/:id", handlers.CategorizeBookmark)
		}

		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		fmt.Printf("Server starting on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited")
}
