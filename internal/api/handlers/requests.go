package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services"
)

type RequestsHandler struct {
	requests services.RequestService
}

func NewRequestsHandler(requests services.RequestService) *RequestsHandler {
	return &RequestsHandler{requests: requests}
}

func (h *RequestsHandler) CreateRequest(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.CreateHelpRequestInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	request, err := h.requests.Create(c.Request.Context(), user.ID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, request)
}

func (h *RequestsHandler) ListRequests(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var filter models.RequestListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	requests, err := h.requests.List(c.Request.Context(), user.ID, filter)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, requests)
}

func (h *RequestsHandler) GetRequest(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	requestID := c.Param("requestId")
	request, err := h.requests.Get(c.Request.Context(), user.ID, requestID)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, request)
}

func (h *RequestsHandler) CancelRequest(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.CancelRequestInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	requestID := c.Param("requestId")
	request, err := h.requests.Cancel(c.Request.Context(), user.ID, requestID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, request)
}

func (h *RequestsHandler) RateHelper(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.RateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	requestID := c.Param("requestId")
	if err := h.requests.RateHelper(c.Request.Context(), user.ID, requestID, payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
