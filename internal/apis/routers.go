package apis

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/handler"
	"github.com/pri-laxmi/MailFlow/internal/middlewares"
	"github.com/pri-laxmi/MailFlow/internal/queue"
	"github.com/pri-laxmi/MailFlow/internal/repository"
	"github.com/pri-laxmi/MailFlow/internal/service"
	"github.com/pri-laxmi/MailFlow/internal/utils"
)

func RegisterRoutes(router *gin.Engine, cfg *configs.Config, db *gorm.DB, jwtManager *utils.JWTManager, jobQueue *queue.Queue) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	contactRepo := repository.NewContactRepository(db)
	contactService := service.NewContactService(contactRepo)
	contactHandler := handler.NewContactHandler(contactService)

	templateRepo := repository.NewTemplateRepository(db)
	templateService := service.NewTemplateService(templateRepo)
	TemplateHandler := handler.NewTemplateHandler(templateService)

	campaignRepo := repository.NewCampaignRepository(db)
	jobRepo := repository.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo, jobQueue)
	campaignService := service.NewCampaignService(
		campaignRepo,
		contactRepo,
		jobRepo,
		jobService,
		jobQueue,
	)
	campaignHandler := handler.NewCampaignHandler(
		campaignService,
	)

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
	protected := router.Group("/")
	protected.Use(middlewares.AuthMIddleware(jwtManager))

	protected.POST("/contacts", contactHandler.Create)
	protected.GET("/contacts", contactHandler.GetAll)
	protected.GET("/contacts/:id", contactHandler.GetByID)
	protected.PATCH("/contacts/:id", contactHandler.Update)
	protected.DELETE("/contacts/:id", contactHandler.Delete)

	protected.POST("/templates", TemplateHandler.Create)
	protected.GET("/templates", TemplateHandler.GetAll)
	protected.GET("/templates/:id", TemplateHandler.GetByID)
	protected.PATCH("/templates/:id", TemplateHandler.Update)
	protected.DELETE("/templates/:id", TemplateHandler.Delete)

	protected.POST("/campaigns", campaignHandler.Create)
	protected.GET("/campaigns", campaignHandler.GetAll)
	protected.GET("/campaigns/:id", campaignHandler.GetByID)
	protected.PATCH("/campaigns/:id", campaignHandler.Update)
	protected.DELETE("/campaigns/:id", campaignHandler.Delete)
}
