package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/api/middleware"
	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
)

func currentUser(c *gin.Context) (*models.User, bool) {
	value, exists := c.Get(middleware.ContextUserKey)
	if !exists {
		return nil, false
	}
	user, ok := value.(*models.User)
	return user, ok
}
