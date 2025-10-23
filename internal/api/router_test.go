package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/MuhibNayem/community-helper-app/internal/api"
	"github.com/MuhibNayem/community-helper-app/internal/api/handlers"
	"github.com/MuhibNayem/community-helper-app/internal/api/middleware"
	"github.com/MuhibNayem/community-helper-app/internal/config"
	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
	"github.com/MuhibNayem/community-helper-app/internal/domain/services/memory"
)

func setupRouter(t *testing.T) (*gin.Engine, *memory.Store) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{HTTPPort: "8080", Env: "test"}
	store := memory.NewStore()

	handlerSet := api.HandlerSet{
		Auth:     handlers.NewAuthHandler(store),
		Users:    handlers.NewUsersHandler(store),
		Requests: handlers.NewRequestsHandler(store),
		Matches:  handlers.NewMatchesHandler(store),
		Health:   handlers.NewHealthHandler(cfg.Env),
	}

	router := api.NewRouter(cfg, handlerSet, middleware.NewAuthMiddleware(store))
	return router, store
}

func doRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var payload *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		payload = bytes.NewReader(data)
	} else {
		payload = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func decodeBody[T any](t *testing.T, rr *httptest.ResponseRecorder, dest *T) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), dest); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

func authenticate(t *testing.T, router *gin.Engine, phone string) (string, models.User) {
	t.Helper()

	resp := doRequest(t, router, http.MethodPost, "/v1/auth/otp/request", gin.H{
		"phone": phone,
	}, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("request otp status = %d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/auth/otp/verify", gin.H{
		"phone":    phone,
		"otp":      "123456",
		"deviceId": "test-device",
	}, "")

	if resp.Code != http.StatusOK {
		t.Fatalf("verify otp status = %d body=%s", resp.Code, resp.Body.String())
	}

	var session models.Session
	decodeBody(t, resp, &session)
	return session.Token, session.User
}

func TestAuthFlow(t *testing.T) {
	router, _ := setupRouter(t)

	resp := doRequest(t, router, http.MethodPost, "/v1/auth/otp/request", gin.H{
		"phone": "+8801000000001",
	}, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("request otp status = %d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/auth/otp/verify", gin.H{
		"phone":    "+8801000000001",
		"otp":      "123456",
		"deviceId": "device-1",
	}, "")

	if resp.Code != http.StatusOK {
		t.Fatalf("verify otp status = %d body=%s", resp.Code, resp.Body.String())
	}

	var session models.Session
	decodeBody(t, resp, &session)
	if session.Token == "" || session.RefreshToken == "" {
		t.Fatalf("expected tokens in session response: %+v", session)
	}
}

func TestProtectedEndpointRequiresAuth(t *testing.T) {
	router, _ := setupRouter(t)
	resp := doRequest(t, router, http.MethodGet, "/v1/requests", nil, "")
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestUsersEndpoints(t *testing.T) {
	router, _ := setupRouter(t)
	token, sessionUser := authenticate(t, router, "+8801000000002")

	resp := doRequest(t, router, http.MethodPatch, "/v1/users/me", gin.H{
		"name":     "Updated User",
		"language": "en",
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("update profile status=%d body=%s", resp.Code, resp.Body.String())
	}

	var updated models.User
	decodeBody(t, resp, &updated)
	if updated.Name != "Updated User" || updated.Language != "en" {
		t.Fatalf("unexpected updated user: %+v", updated)
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/helpers/me/toggle", gin.H{
		"optedIn": false,
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("toggle helper status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPut, "/v1/helpers/me/skills", gin.H{
		"skills": []string{"MEDICAL_FIRST_AID", "MECHANICAL"},
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("update skills status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPut, "/v1/helpers/me/availability", gin.H{
		"weekly": []gin.H{
			{"day": "MONDAY", "start": "09:00", "end": "17:00"},
			{"day": "TUESDAY", "start": "10:00", "end": "18:00"},
		},
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("manage availability status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/helpers/me/kyc", gin.H{
		"documentType": "NID",
		"fileUrl":      "https://example.com/nid.pdf",
	}, token)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("upload kyc status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodGet, "/v1/users/me", nil, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("get current user status=%d body=%s", resp.Code, resp.Body.String())
	}

	var current models.User
	decodeBody(t, resp, &current)
	if current.ID != sessionUser.ID {
		t.Fatalf("expected user id %s got %s", sessionUser.ID, current.ID)
	}
}

func TestRequestLifecycle(t *testing.T) {
	router, store := setupRouter(t)
	token, user := authenticate(t, router, "+8801000000003")

	// Create request A and cancel it
	reqBody := gin.H{
		"type":        "URGENT",
		"category":    "MEDICAL_FIRST_AID",
		"description": "Need immediate assistance",
		"location": gin.H{
			"lat":     23.78,
			"lng":     90.36,
			"address": "Dhaka",
		},
	}
	resp := doRequest(t, router, http.MethodPost, "/v1/requests", reqBody, token)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create request status=%d body=%s", resp.Code, resp.Body.String())
	}
	var requestA models.HelpRequest
	decodeBody(t, resp, &requestA)

	resp = doRequest(t, router, http.MethodPost, "/v1/requests/"+requestA.ID+"/cancel", gin.H{
		"reason": "HELPER_NOT_NEEDED",
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("cancel request status=%d body=%s", resp.Code, resp.Body.String())
	}

	// Create request B for full lifecycle
	reqBodyB := gin.H{
		"type":        "PLANNED",
		"category":    "TUTORING",
		"description": "Schedule future session",
		"location": gin.H{
			"lat":     23.8,
			"lng":     90.4,
			"address": "Dhaka",
		},
		"scheduledFor": time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/requests", reqBodyB, token)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create requestB status=%d body=%s", resp.Code, resp.Body.String())
	}
	var requestB models.HelpRequest
	decodeBody(t, resp, &requestB)

	match := store.SeedMatch(user.ID, requestB.ID)

	resp = doRequest(t, router, http.MethodGet, "/v1/matches", nil, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("list matches status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/matches/"+match.ID+"/accept", nil, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("accept match status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/matches/"+match.ID+"/status", gin.H{
		"status": "COMPLETED",
	}, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("update match status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodPost, "/v1/requests/"+requestB.ID+"/rate", gin.H{
		"rating": 5,
	}, token)
	if resp.Code != http.StatusCreated {
		t.Fatalf("rate helper status=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = doRequest(t, router, http.MethodGet, "/v1/requests/"+requestB.ID, nil, token)
	if resp.Code != http.StatusOK {
		t.Fatalf("get request status=%d body=%s", resp.Code, resp.Body.String())
	}
	var fetched models.HelpRequest
	decodeBody(t, resp, &fetched)
	if fetched.Status != "COMPLETED" {
		t.Fatalf("expected completed status got %s", fetched.Status)
	}
}

