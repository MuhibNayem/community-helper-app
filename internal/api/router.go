package api

import (
	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/api/handlers"
	"github.com/MuhibNayem/community-helper-app/internal/config"
)

type HandlerSet struct {
	Auth     *handlers.AuthHandler
	Users    *handlers.UsersHandler
	Requests *handlers.RequestsHandler
	Matches  *handlers.MatchesHandler
	Health   *handlers.HealthHandler
}

func NewRouter(cfg *config.Config, handlers HandlerSet, authMiddleware gin.HandlerFunc) *gin.Engine {
	if cfg.Env == "production" || cfg.Env == "stage" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.SetTrustedProxies(nil)
	engine.Use(gin.Recovery())

	engine.GET("/healthz", handlers.Health.Check)

	v1 := engine.Group("/v1")

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/otp/request", handlers.Auth.RequestOTP)
		authGroup.POST("/otp/verify", handlers.Auth.VerifyOTP)
		authGroup.POST("/refresh", handlers.Auth.RefreshSession)
		authGroup.POST("/logout", handlers.Auth.Logout)
	}

	protected := v1.Group("")
	protected.Use(authMiddleware)

	protected.GET("/users/me", handlers.Users.GetCurrentUser)
	protected.PATCH("/users/me", handlers.Users.UpdateProfile)
	protected.POST("/helpers/me/toggle", handlers.Users.ToggleHelper)
	protected.PUT("/helpers/me/skills", handlers.Users.UpdateSkills)
	protected.PUT("/helpers/me/availability", handlers.Users.ManageAvailability)
	protected.POST("/helpers/me/kyc", handlers.Users.UploadKYC)

	protected.POST("/requests", handlers.Requests.CreateRequest)
	protected.GET("/requests", handlers.Requests.ListRequests)
	protected.GET("/requests/:requestId", handlers.Requests.GetRequest)
	protected.POST("/requests/:requestId/cancel", handlers.Requests.CancelRequest)
	protected.POST("/requests/:requestId/rate", handlers.Requests.RateHelper)

	protected.GET("/matches", handlers.Matches.ListInvitations)
	protected.POST("/matches/:matchId/accept", handlers.Matches.AcceptInvitation)
	protected.POST("/matches/:matchId/decline", handlers.Matches.DeclineInvitation)
	protected.POST("/matches/:matchId/status", handlers.Matches.UpdateStatus)

	return engine
}
