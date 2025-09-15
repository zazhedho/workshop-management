package router

import (
	"net/http"
	bookingHandler "workshop-management/internal/handlers/http/booking"
	serviceHandler "workshop-management/internal/handlers/http/service"
	userHandler "workshop-management/internal/handlers/http/user"
	vehicleHandler "workshop-management/internal/handlers/http/vehicle"
	authRepo "workshop-management/internal/repositories/auth"
	bookingRepo "workshop-management/internal/repositories/booking"
	serviceRepo "workshop-management/internal/repositories/service"
	userRepo "workshop-management/internal/repositories/user"
	vehicleRepo "workshop-management/internal/repositories/vehicle"
	bookingSvc "workshop-management/internal/services/booking"
	serviceSvc "workshop-management/internal/services/service"
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
			userPriv.PUT("", h.Update)
			userPriv.DELETE("", h.Delete)
		}
	}

	r.App.GET("/api/users", mdw.AuthMiddleware(), mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCashier), h.GetAllUsers)
}

func (r *Routes) VehicleRoutes() {
	repo := vehicleRepo.NewVehicleRepo(r.DB)
	uc := vehicleSvc.NewVehicleService(repo)
	h := vehicleHandler.NewVehicleHandler(uc)
	mdw := middlewares.NewMiddleware(authRepo.NewBlacklistRepo(r.DB))

	vehicle := r.App.Group("/api/vehicle").Use(mdw.AuthMiddleware())
	{
		vehicle.POST("", h.Create)
		vehicle.GET("/:id", h.GetById)
		vehicle.GET("/list", h.Fetch)
		vehicle.PUT("/:id", mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCustomer), h.Update)
		vehicle.DELETE("/:id", mdw.RoleMiddleware(utils.RoleAdmin, utils.RoleCustomer), h.Delete)
	}

}

func (r *Routes) ServiceRoutes() {
	repo := serviceRepo.NewServiceRepo(r.DB)
	uc := serviceSvc.NewSrvService(repo)
	h := serviceHandler.NewServiceHandler(uc)
	mdw := middlewares.NewMiddleware(authRepo.NewBlacklistRepo(r.DB))

	svc := r.App.Group("/api/services")
	{
		svc.GET("", h.Fetch)
		svc.GET("/:id", h.GetById)

		svcPriv := svc.Group("").Use(mdw.AuthMiddleware(), mdw.RoleMiddleware(utils.RoleAdmin))
		{
			svcPriv.POST("", h.Create)
			svcPriv.PUT("/:id", h.Update)
			svcPriv.DELETE("/:id", h.Delete)
		}
	}
}

func (r *Routes) BookingRoutes() {
	repo := bookingRepo.NewBookingRepo(r.DB)
	uc := bookingSvc.NewServiceBooking(repo)
	h := bookingHandler.NewBookingHandler(uc)
	mdw := middlewares.NewMiddleware(authRepo.NewBlacklistRepo(r.DB))

	booking := r.App.Group("/api/bookings").Use(mdw.AuthMiddleware())
	{
		booking.POST("", h.Create)
		booking.GET("/:id", h.GetBookingById)
		booking.GET("", h.Fetch)
	}
}
