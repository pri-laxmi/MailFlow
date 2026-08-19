package apis

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/handler"
	"github.com/pri-laxmi/MailFlow/internal/repository"
	"github.com/pri-laxmi/MailFlow/internal/service"
	"github.com/pri-laxmi/MailFlow/internal/utils"
)

func RegisterRoutes(router *gin.Engine, cfg *configs.Config, db *gorm.DB, jwtManager *utils.JWTManager) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)

	authHandler := handler.NewAuthHandler(authService)
	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	router.GET("/ready", func(c *gin.Context) {
		err := db.Exec("SELECT 1").Error
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
}
