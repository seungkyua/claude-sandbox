package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/ktc-plugin-hub/backend/internal/handler"
	"github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/ktc-plugin-hub/backend/internal/service"
	"gorm.io/gorm"
)

// SetupRouter 는 전체 라우트를 등록하고 Gin 엔진을 반환한다
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 미들웨어 적용
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.JSONErrorHandler())

	// Repository 초기화
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	pluginRepo := repository.NewPluginRepository(db)
	versionRepo := repository.NewVersionRepository(db)
	installRepo := repository.NewInstallationRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Service 초기화
	authSvc := service.NewAuthService(userRepo, &cfg.JWT)
	pluginSvc := service.NewPluginService(pluginRepo, categoryRepo)
	versionSvc := service.NewVersionService(versionRepo, pluginRepo)
	installSvc := service.NewInstallationService(installRepo, pluginRepo, versionRepo)
	reviewSvc := service.NewReviewService(reviewRepo, pluginRepo)
	adminSvc := service.NewAdminService(pluginRepo)

	// Handler 초기화
	authHandler := handler.NewAuthHandler(authSvc)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	pluginHandler := handler.NewPluginHandler(pluginSvc)
	versionHandler := handler.NewVersionHandler(versionSvc)
	installHandler := handler.NewInstallationHandler(installSvc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)

	// API v1 그룹
	v1 := r.Group("/api/v1")
	{
		// 공개 라우트 (인증 불필요)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		v1.GET("/categories", categoryHandler.GetAll)
		v1.GET("/plugins", pluginHandler.GetList)
		v1.GET("/plugins/:id", pluginHandler.GetByID)
		v1.GET("/plugins/:id/reviews", reviewHandler.GetReviews)

		// 인증 필요 라우트
		authenticated := v1.Group("", middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			authenticated.GET("/me", authHandler.Me)

			authenticated.POST("/plugins", pluginHandler.Create)
			authenticated.PUT("/plugins/:id", pluginHandler.Update)
			authenticated.DELETE("/plugins/:id", pluginHandler.Delete)

			authenticated.POST("/plugins/:id/versions", versionHandler.CreateVersion)
			authenticated.GET("/plugins/:id/versions/:versionId/download", versionHandler.Download)

			authenticated.POST("/plugins/:id/install", installHandler.Install)
			authenticated.DELETE("/plugins/:id/install", installHandler.Uninstall)
			authenticated.PATCH("/plugins/:id/install", installHandler.ToggleActive)
			authenticated.GET("/me/installations", installHandler.GetMyInstallations)

			authenticated.POST("/plugins/:id/reviews", reviewHandler.CreateReview)
			authenticated.PUT("/plugins/:id/reviews/:reviewId", reviewHandler.UpdateReview)
			authenticated.DELETE("/plugins/:id/reviews/:reviewId", reviewHandler.DeleteReview)
		}

		// 관리자 전용 라우트
		admin := v1.Group("/admin", middleware.AuthMiddleware(cfg.JWT.Secret), middleware.AdminMiddleware())
		{
			admin.GET("/plugins/pending", adminHandler.GetPendingPlugins)
			admin.PATCH("/plugins/:id/approve", adminHandler.ApprovePlugin)
			admin.PATCH("/plugins/:id/reject", adminHandler.RejectPlugin)
			admin.PATCH("/plugins/:id/hide", adminHandler.HidePlugin)
		}
	}

	return r
}
