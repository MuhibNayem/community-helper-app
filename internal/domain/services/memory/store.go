package memory

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
)

var (
	errUserNotFound    = errors.New("user not found")
	errRequestNotFound = errors.New("request not found")
	errMatchNotFound   = errors.New("match not found")
	errUnauthorized    = errors.New("unauthorized")
)

type Store struct {
	mu sync.RWMutex

	now func() time.Time

	users          map[string]*models.User
	helperProfiles map[string]*models.HelperProfile
	kycDocuments   map[string]*models.KYCDocument
	requests       map[string]*models.HelpRequest
	matches        map[string]*models.MatchSession

	otps          map[string]string
	sessions      map[string]*models.Session
	refreshTokens map[string]string

	nextRequestID int
	nextMatchID   int
}

func NewStore() *Store {
	return &Store{
		now:            time.Now,
		users:          make(map[string]*models.User),
		helperProfiles: make(map[string]*models.HelperProfile),
		kycDocuments:   make(map[string]*models.KYCDocument),
		requests:       make(map[string]*models.HelpRequest),
		matches:        make(map[string]*models.MatchSession),
		otps:           make(map[string]string),
		sessions:       make(map[string]*models.Session),
		refreshTokens:  make(map[string]string),
		nextRequestID:  1,
		nextMatchID:    1,
	}
}

func (s *Store) withNow(now func() time.Time) *Store {
	s.now = now
	return s
}

// AuthService implementation

func (s *Store) RequestOTP(_ context.Context, req models.OTPRequest) (models.OTPRequestResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.otps[req.Phone] = "123456"

	return models.OTPRequestResponse{
		ExpiresIn:         120,
		AttemptsRemaining: 3,
	}, nil
}

func (s *Store) VerifyOTP(_ context.Context, req models.OTPVerification) (models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expected := s.otps[req.Phone]
	if expected == "" || expected != req.OTP {
		return models.Session{}, fmt.Errorf("invalid otp")
	}

	user := s.ensureUser(req.Phone)
	session := s.createSession(user.ID)

	return *session, nil
}

func (s *Store) RefreshSession(_ context.Context, refreshToken string) (models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID, ok := s.refreshTokens[refreshToken]
	if !ok {
		return models.Session{}, fmt.Errorf("invalid refresh token")
	}

	session := s.createSession(userID)
	return *session, nil
}

func (s *Store) Logout(_ context.Context, deviceID string) error {
	_ = deviceID
	return nil
}

func (s *Store) ParseToken(_ context.Context, token string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[token]
	if !ok {
		return nil, errUnauthorized
	}

	user, ok := s.users[session.User.ID]
	if !ok {
		return nil, errUserNotFound
	}

	copied := *user
	return &copied, nil
}

// UserService implementation

func (s *Store) GetCurrentUser(_ context.Context, userID string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, errUserNotFound
	}

	copied := *user
	return &copied, nil
}

func (s *Store) UpdateProfile(_ context.Context, userID string, update models.UserUpdate) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, errUserNotFound
	}

	if update.Name != nil {
		user.Name = *update.Name
	}
	if update.Language != nil {
		user.Language = strings.ToLower(*update.Language)
	}
	if update.PhotoURL != nil {
		user.PhotoURL = *update.PhotoURL
	}
	user.UpdatedAt = s.now()

	copied := *user
	return &copied, nil
}

func (s *Store) ToggleHelper(_ context.Context, userID string, toggle models.HelperToggle) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, errUserNotFound
	}

	user.IsHelper = toggle.OptedIn
	if toggle.OptedIn {
		user.HelperStatus = "ACTIVE"
	} else {
		user.HelperStatus = "INACTIVE"
	}
	user.UpdatedAt = s.now()

	profile := s.ensureHelperProfile(userID)
	profile.OptedIn = toggle.OptedIn
	profile.UpdatedAt = s.now()

	copied := *user
	return &copied, nil
}

func (s *Store) UpdateSkills(_ context.Context, userID string, update models.SkillsUpdate) (*models.HelperProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	profile := s.ensureHelperProfile(userID)
	if len(update.Skills) == 0 {
		return nil, fmt.Errorf("skills required")
	}

	profile.Skills = append([]string{}, update.Skills...)
	profile.UpdatedAt = s.now()

	copyProfile := *profile
	return &copyProfile, nil
}

func (s *Store) ManageAvailability(_ context.Context, userID string, update models.AvailabilityUpdate) (*models.HelperProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(update.Weekly) == 0 {
		return nil, fmt.Errorf("weekly availability required")
	}

	profile := s.ensureHelperProfile(userID)

	profile.Availability = models.Availability{
		Weekly:     make([]models.AvailabilitySlot, len(update.Weekly)),
		Exceptions: make([]models.AvailabilityException, len(update.Exceptions)),
	}

	for i, slot := range update.Weekly {
		profile.Availability.Weekly[i] = models.AvailabilitySlot{
			Day:   slot.Day,
			Start: slot.Start,
			End:   slot.End,
		}
	}

	for i, ex := range update.Exceptions {
		exception := models.AvailabilityException{
			Date:  ex.Date,
			Slots: make([]models.ExceptionSlot, len(ex.Slots)),
		}
		for j, slot := range ex.Slots {
			exception.Slots[j] = models.ExceptionSlot{
				Start: slot.Start,
				End:   slot.End,
			}
		}
		profile.Availability.Exceptions[i] = exception
	}

	profile.UpdatedAt = s.now()

	copyProfile := *profile
	return &copyProfile, nil
}

func (s *Store) UploadKYC(_ context.Context, userID string, upload models.KYCDocumentUpload) (*models.KYCDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if upload.DocumentType == "" || upload.FileURL == "" {
		return nil, fmt.Errorf("document type and fileUrl required")
	}

	id := fmt.Sprintf("kyc-%d", len(s.kycDocuments)+1)
	doc := &models.KYCDocument{
		ID:           id,
		UserID:       userID,
		DocumentType: upload.DocumentType,
		FileURL:      upload.FileURL,
		Status:       "PENDING",
		SubmittedAt:  s.now(),
	}

	s.kycDocuments[id] = doc

	copyDoc := *doc
	return &copyDoc, nil
}

// RequestService implementation

func (s *Store) Create(_ context.Context, userID string, input models.CreateHelpRequestInput) (*models.HelpRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("req-%d", s.nextRequestID)
	s.nextRequestID++

	now := s.now()
	request := &models.HelpRequest{
		ID:          id,
		RequesterID: userID,
		Type:        input.Type,
		Status:      "SUBMITTED",
		Category:    input.Category,
		Description: input.Description,
		Attachments: append([]string{}, input.Attachments...),
		Location: models.RequestLocation{
			Latitude:  input.Location.Latitude,
			Longitude: input.Location.Longitude,
			Address:   input.Location.Address,
			PlaceID:   input.Location.PlaceID,
			Accuracy:  input.Location.Accuracy,
		},
		ScheduledFor: input.ScheduledFor,
		CreatedAt:    now,
		UpdatedAt:    now,
		SLA: models.SLAWindows{
			MatchDeadline:      now.Add(15 * time.Minute),
			CompletionDeadline: now.Add(6 * time.Hour),
		},
		Pricing: models.Pricing{
			EstimatedAmount: 500,
			Currency:        "BDT",
			PlatformFee:     50,
		},
	}

	s.requests[id] = request

	copyRequest := *request
	return &copyRequest, nil
}

func (s *Store) List(_ context.Context, userID string, filter models.RequestListFilter) ([]models.HelpRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.HelpRequest
	for _, req := range s.requests {
		if req.RequesterID != userID {
			continue
		}
		if filter.Status != "" && req.Status != filter.Status {
			continue
		}
		results = append(results, *req)
	}

	limit := filter.Limit
	if limit <= 0 {
		return results, nil
	}

	offset := filter.Offset
	if offset >= len(results) {
		return []models.HelpRequest{}, nil
	}

	end := offset + limit
	if end > len(results) {
		end = len(results)
	}

	return results[offset:end], nil
}

func (s *Store) Get(_ context.Context, userID, requestID string) (*models.HelpRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	req, ok := s.requests[requestID]
	if !ok || req.RequesterID != userID {
		return nil, errRequestNotFound
	}

	copyReq := *req
	return &copyReq, nil
}

func (s *Store) Cancel(_ context.Context, userID, requestID string, input models.CancelRequestInput) (*models.HelpRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := s.requests[requestID]
	if !ok || req.RequesterID != userID {
		return nil, errRequestNotFound
	}

	if req.Status == "COMPLETED" {
		return nil, fmt.Errorf("cannot cancel completed request")
	}

	now := s.now()
	req.Status = "CANCELLED"
	req.Cancellation = &models.Cancellation{
		Reason:         input.Reason,
		Initiator:      "SEEKER",
		Timestamp:      now,
		PenaltyApplied: false,
	}
	req.UpdatedAt = now

	copyReq := *req
	return &copyReq, nil
}

func (s *Store) RateHelper(_ context.Context, userID, requestID string, rating models.RateRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := s.requests[requestID]
	if !ok || req.RequesterID != userID {
		return errRequestNotFound
	}

	if req.Status != "COMPLETED" {
		return fmt.Errorf("request not completed")
	}

	// Pretend to update helper stats.
	return nil
}

// MatchService implementation

func (s *Store) ListInvitations(_ context.Context, helperID string, status string) ([]models.MatchSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []models.MatchSession
	for _, match := range s.matches {
		if match.HelperID != helperID {
			continue
		}
		if status != "" && match.Status != status {
			continue
		}
		matches = append(matches, *match)
	}
	return matches, nil
}

func (s *Store) Accept(_ context.Context, helperID, matchID string) (*models.MatchSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	match, ok := s.matches[matchID]
	if !ok || match.HelperID != helperID {
		return nil, errMatchNotFound
	}

	now := s.now()
	match.Status = "ACCEPTED"
	match.AcceptedAt = &now
	match.RespondedAt = &now

	if req, ok := s.requests[match.RequestID]; ok {
		req.Status = "ACCEPTED"
		req.UpdatedAt = now
	}

	copyMatch := *match
	return &copyMatch, nil
}

func (s *Store) Decline(_ context.Context, helperID, matchID string, input models.DeclineMatchInput) (*models.MatchSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	match, ok := s.matches[matchID]
	if !ok || match.HelperID != helperID {
		return nil, errMatchNotFound
	}

	now := s.now()
	match.Status = "DECLINED"
	match.DeclineReason = input.Reason
	match.RespondedAt = &now

	copyMatch := *match
	return &copyMatch, nil
}

func (s *Store) UpdateStatus(_ context.Context, helperID, matchID string, input models.MatchStatusUpdate) (*models.MatchSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	match, ok := s.matches[matchID]
	if !ok || match.HelperID != helperID {
		return nil, errMatchNotFound
	}

	now := s.now()
	match.Status = input.Status

	switch input.Status {
	case "EN_ROUTE":
		match.ArrivedAt = nil
	case "ARRIVED":
		match.ArrivedAt = &now
	case "COMPLETED":
		match.CompletedAt = &now
		if req, ok := s.requests[match.RequestID]; ok {
			req.Status = "COMPLETED"
			req.UpdatedAt = now
		}
	}

	copyMatch := *match
	return &copyMatch, nil
}

// Helpers

func (s *Store) ensureUser(phone string) *models.User {
	for _, user := range s.users {
		if user.Phone == phone {
			return user
		}
	}

	id := fmt.Sprintf("user-%d", len(s.users)+1)
	now := s.now()
	user := &models.User{
		ID:        id,
		Phone:     phone,
		Name:      "New User",
		Language:  "bn",
		IsHelper:  true,
		CreatedAt: now,
		UpdatedAt: now,
		NotificationPrefs: models.NotificationPrefs{
			QuietHours: models.QuietHours{
				Start: "22:00",
				End:   "07:00",
			},
			UrgentSMS: true,
		},
		HelperStatus: "ACTIVE",
		KYCStatus:    "PENDING",
	}

	s.users[id] = user
	s.ensureHelperProfile(id)

	return user
}

func (s *Store) ensureHelperProfile(userID string) *models.HelperProfile {
	profile, ok := s.helperProfiles[userID]
	if ok {
		return profile
	}

	now := s.now()
	profile = &models.HelperProfile{
		UserID:    userID,
		Skills:    []string{"GENERAL_HELP"},
		OptedIn:   true,
		Rating:    5,
		Badges:    []string{},
		UpdatedAt: now,
		Availability: models.Availability{
			Weekly: []models.AvailabilitySlot{
				{Day: "MONDAY", Start: "09:00", End: "17:00"},
			},
		},
	}

	s.helperProfiles[userID] = profile
	return profile
}

func (s *Store) createSession(userID string) *models.Session {
	token := fmt.Sprintf("token-%s-%d", userID, time.Now().UnixNano())
	refresh := fmt.Sprintf("refresh-%s-%d", userID, time.Now().UnixNano())

	user := s.users[userID]

	session := &models.Session{
		Token:        token,
		RefreshToken: refresh,
		User:         *user,
		ExpiresAt:    s.now().Add(24 * time.Hour),
	}

	s.sessions[token] = session
	s.refreshTokens[refresh] = userID

	return session
}

func (s *Store) SeedMatch(helperID, requestID string) *models.MatchSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("match-%d", s.nextMatchID)
	s.nextMatchID++

	now := s.now()
	match := &models.MatchSession{
		ID:        id,
		RequestID: requestID,
		HelperID:  helperID,
		Status:    "INVITED",
		InvitedAt: now,
	}

	s.matches[id] = match
	return match
}
