package router

import (
	"net/http"
	userHandler "workshop-management/internal/handlers/http/user"
	vehicleHandler "workshop-management/internal/handlers/http/vehicle"
	authRepo "workshop-management/internal/repositories/auth"
	userRepo "workshop-management/internal/repositories/user"
	vehicleRepo "workshop-management/internal/repositories/vehicle"
	userSvc "workshop-management/internal/services/user"
	vehicleSvc "workshop-management/internal/services/vehicle"
	"workshop-management/middlewares"
	"workshop-management/pkg/logger"
	"workshop-management/utils"

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

	app.Use(middlewares.CORS())
	app.Use(gin.CustomRecovery(middlewares.ErrorHandler))
	app.Use(middlewares.SetContextId())

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
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	repo := userRepo.NewUserRepo(r.DB)
	uc := userSvc.NewUserService(repo, blacklistRepo)
	h := userHandler.NewUserHandler(uc)
	mdw := middlewares.NewMiddleware(blacklistRepo)

	user := r.App.Group("/api/user")
	{
		user.POST("/register", h.Register)
		user.POST("/login", h.Login)

		userPriv := user.Group("").Use(mdw.AuthMiddleware())
		{
			userPriv.POST("/logout", h.Logout)
			userPriv.GET("", h.GetUserByAuth)
			userPriv.GET("/:id", mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCashier), h.GetUserById)
			userPriv.PUT("", h.UpdateUser)
			userPriv.DELETE("/:id", h.DeleteUser)
		}
	}

	r.App.GET("/api/users", mdw.AuthMiddleware(), mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCashier), h.GetAllUsers)
}

func (r *Routes) VehicleRoutes() {
	repo := vehicleRepo.NewVehicleRepo(r.DB)
	uc := vehicleSvc.NewVehicleService(repo)
	h := vehicleHandler.NewVehicleHandler(uc)

	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo)

	vehicle := r.App.Group("/api/vehicle").Use(mdw.AuthMiddleware())
	{
		vehicle.POST("", h.CreateVehicle)
		vehicle.GET("/:id", h.GetVehicle)
		vehicle.GET("/list", h.FetchVehicles)
		//vehicle.PUT("/:id", mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCustomer), h.UpdateVehicle)
		//vehicle.DELETE("/:id", h.DeleteVehicle)
	}

}
