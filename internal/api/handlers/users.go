package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services"
)

type UsersHandler struct {
	users services.UserService
}

func NewUsersHandler(users services.UserService) *UsersHandler {
	return &UsersHandler{users: users}
}

func (h *UsersHandler) GetCurrentUser(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	result, err := h.users.GetCurrentUser(c.Request.Context(), user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, result)
}

func (h *UsersHandler) UpdateProfile(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.UserUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.users.UpdateProfile(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, result)
}

func (h *UsersHandler) ToggleHelper(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.HelperToggle
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.users.ToggleHelper(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, result)
}

func (h *UsersHandler) UpdateSkills(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.SkillsUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.users.UpdateSkills(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, result)
}

func (h *UsersHandler) ManageAvailability(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.AvailabilityUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.users.ManageAvailability(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, result)
}

func (h *UsersHandler) UploadKYC(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.KYCDocumentUpload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.users.UploadKYC(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusAccepted, result)
}
