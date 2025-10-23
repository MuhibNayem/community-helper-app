package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services"
)

type MatchesHandler struct {
	matches services.MatchService
}

func NewMatchesHandler(matches services.MatchService) *MatchesHandler {
	return &MatchesHandler{matches: matches}
}

func (h *MatchesHandler) ListInvitations(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	status := c.Query("status")
	results, err := h.matches.ListInvitations(c.Request.Context(), user.ID, status)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, results)
}

func (h *MatchesHandler) AcceptInvitation(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	matchID := c.Param("matchId")
	match, err := h.matches.Accept(c.Request.Context(), user.ID, matchID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, match)
}

func (h *MatchesHandler) DeclineInvitation(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.DeclineMatchInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	matchID := c.Param("matchId")
	match, err := h.matches.Decline(c.Request.Context(), user.ID, matchID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, match)
}

func (h *MatchesHandler) UpdateStatus(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload models.MatchStatusUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	matchID := c.Param("matchId")
	match, err := h.matches.UpdateStatus(c.Request.Context(), user.ID, matchID, payload)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, match)
}
