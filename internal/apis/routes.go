package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/utils"
	"gorm.io/gorm"
)

func SetRoutes(cfg *configs.Config, db *gorm.DB, jwtManager *utils.JWTManager) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, cfg, db, jwtManager)
	return router
}
