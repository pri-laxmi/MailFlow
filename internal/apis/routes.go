package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/pri-laxmi/MailFlow/configs"
	"gorm.io/gorm"
)

func SetRoutes(cfg *configs.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, cfg, db)
	return router
}
