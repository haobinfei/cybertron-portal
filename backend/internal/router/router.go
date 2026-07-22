package router

import (
	"cybertron-portal/internal/handler"
	"cybertron-portal/internal/middleware"
	"cybertron-portal/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(
	authService *service.AuthService,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.MetricsMiddleware())

	// Prometheus metrics endpoint (no auth required)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			authorized.POST("/auth/logout", authHandler.Logout)

			user := authorized.Group("/user")
			{
				user.GET("/me", userHandler.GetMe)
			}

			admin := authorized.Group("/users")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("", userHandler.List)
				admin.POST("", userHandler.Create)
				admin.PUT("/:id", userHandler.Update)
				admin.DELETE("/:id", userHandler.Delete)
			}
		}
	}

	return r
}
