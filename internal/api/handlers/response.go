package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, payload)
}

func writeError(c *gin.Context, status int, msg string) {
	writeJSON(c, status, errorResponse{Error: msg})
}

func notImplemented(c *gin.Context) {
	writeError(c, http.StatusNotImplemented, "not implemented")
}
