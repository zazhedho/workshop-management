package http

import (
	"net/http"
	"workshop-management/internal/handler/http/users"
	"workshop-management/internal/repository/auth"
	"workshop-management/internal/repository/users"
	"workshop-management/internal/services/users"
	"workshop-management/middleware"
	"workshop-management/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Routes struct {
	App *gin.Engine
	DB  *gorm.DB
}

func NewRoutes() *Routes {
	app := gin.Default()

	app.Use(middleware.CORS())
	app.Use(gin.CustomRecovery(middleware.ErrorHandler))
	app.Use(middleware.SetContextId())

	// health check
	app.GET("/healthcheck", func(ctx *gin.Context) {
		logger.WriteLog(logger.LogLevelDebug, "ClientIP: "+ctx.ClientIP())
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK!!",
		})
	})
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Routes{
		App: app,
	}
}

func (r *Routes) UserRoutes() {
	blacklistRepo := auth.NewBlacklistRepo(r.DB)
	repo := users.NewUserRepo(r.DB)
	uc := service.NewUserService(repo, blacklistRepo)
	h := handler.NewUserHandler(uc)
	mdw := middleware.NewMiddleware(blacklistRepo)

	user := r.App.Group("/api/user")
	{
		user.POST("/register", h.Register)
		user.POST("/login", h.Login)
		user.POST("/logout", mdw.AuthMiddleware(), h.Logout)
	}
}
