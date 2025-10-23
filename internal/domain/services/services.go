package services

import (
	"context"

	"github.com/MuhibNayem/community-helper-app/internal/domain/models"
)

type AuthService interface {
	RequestOTP(ctx context.Context, req models.OTPRequest) (models.OTPRequestResponse, error)
	VerifyOTP(ctx context.Context, req models.OTPVerification) (models.Session, error)
	RefreshSession(ctx context.Context, refreshToken string) (models.Session, error)
	Logout(ctx context.Context, deviceID string) error
	ParseToken(ctx context.Context, token string) (*models.User, error)
}

type UserService interface {
	GetCurrentUser(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, userID string, update models.UserUpdate) (*models.User, error)
	ToggleHelper(ctx context.Context, userID string, toggle models.HelperToggle) (*models.User, error)
	UpdateSkills(ctx context.Context, userID string, update models.SkillsUpdate) (*models.HelperProfile, error)
	ManageAvailability(ctx context.Context, userID string, update models.AvailabilityUpdate) (*models.HelperProfile, error)
	UploadKYC(ctx context.Context, userID string, upload models.KYCDocumentUpload) (*models.KYCDocument, error)
}

type RequestService interface {
	Create(ctx context.Context, userID string, input models.CreateHelpRequestInput) (*models.HelpRequest, error)
	List(ctx context.Context, userID string, filter models.RequestListFilter) ([]models.HelpRequest, error)
	Get(ctx context.Context, userID, requestID string) (*models.HelpRequest, error)
	Cancel(ctx context.Context, userID, requestID string, input models.CancelRequestInput) (*models.HelpRequest, error)
	RateHelper(ctx context.Context, userID, requestID string, rating models.RateRequest) error
}

type MatchService interface {
	ListInvitations(ctx context.Context, helperID string, status string) ([]models.MatchSession, error)
	Accept(ctx context.Context, helperID, matchID string) (*models.MatchSession, error)
	Decline(ctx context.Context, helperID, matchID string, input models.DeclineMatchInput) (*models.MatchSession, error)
	UpdateStatus(ctx context.Context, helperID, matchID string, input models.MatchStatusUpdate) (*models.MatchSession, error)
}
