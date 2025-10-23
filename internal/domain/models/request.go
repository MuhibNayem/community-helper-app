package models

import "time"

type HelpRequest struct {
	ID           string          `json:"id"`
	RequesterID  string          `json:"requesterId"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	Category     string          `json:"category"`
	Description  string          `json:"description"`
	Attachments  []string        `json:"attachments,omitempty"`
	Location     RequestLocation `json:"location"`
	ScheduledFor *time.Time      `json:"scheduledFor,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	ExpiresAt    *time.Time      `json:"expiresAt,omitempty"`
	SLA          SLAWindows      `json:"sla"`
	Pricing      Pricing         `json:"pricing"`
	Cancellation *Cancellation   `json:"cancellation,omitempty"`
}

type RequestLocation struct {
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	PlaceID   string  `json:"placeId,omitempty"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

type SLAWindows struct {
	MatchDeadline      time.Time `json:"matchDeadline"`
	CompletionDeadline time.Time `json:"completionDeadline"`
}

type Pricing struct {
	EstimatedAmount float64 `json:"estimatedAmount"`
	Currency        string  `json:"currency"`
	PlatformFee     float64 `json:"platformFee"`
}

type Cancellation struct {
	Reason         string    `json:"reason"`
	Initiator      string    `json:"initiator"`
	Timestamp      time.Time `json:"timestamp"`
	PenaltyApplied bool      `json:"penaltyApplied"`
}

type RateRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment,omitempty" binding:"omitempty,max=500"`
}

type CreateHelpRequestInput struct {
	Type         string          `json:"type" binding:"required,oneof=URGENT PLANNED"`
	Category     string          `json:"category" binding:"required"`
	Description  string          `json:"description" binding:"required"`
	Location     RequestLocation `json:"location" binding:"required"`
	ScheduledFor *time.Time      `json:"scheduledFor,omitempty" binding:"omitempty"`
	Attachments  []string        `json:"attachments,omitempty"`
}

type RequestListFilter struct {
	Status string `form:"status"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type CancelRequestInput struct {
	Reason string `json:"reason" binding:"required"`
}
