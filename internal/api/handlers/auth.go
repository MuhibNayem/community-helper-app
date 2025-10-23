package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var payload models.OTPRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.RequestOTP(c.Request.Context(), payload)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var payload models.OTPVerification
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.authService.VerifyOTP(c.Request.Context(), payload)
	if err != nil {
		writeError(c, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, session)
}

func (h *AuthHandler) RefreshSession(c *gin.Context) {
	var payload models.RefreshRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.authService.RefreshSession(c.Request.Context(), payload.RefreshToken)
	if err != nil {
		writeError(c, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, session)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var payload models.LogoutRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Logout(c.Request.Context(), payload.DeviceID); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
