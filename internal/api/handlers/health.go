package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	env string
}

func NewHealthHandler(env string) *HealthHandler {
	return &HealthHandler{env: env}
}

func (h *HealthHandler) Check(c *gin.Context) {
	writeJSON(c, http.StatusOK, map[string]string{
		"status": "ok",
		"env":    h.env,
	})
}
