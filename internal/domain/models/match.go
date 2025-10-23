package models

import "time"

type MatchSession struct {
	ID            string        `json:"id"`
	RequestID     string        `json:"requestId"`
	HelperID      string        `json:"helperId"`
	Status        string        `json:"status"`
	InvitedAt     time.Time     `json:"invitedAt"`
	RespondedAt   *time.Time    `json:"respondedAt,omitempty"`
	AcceptedAt    *time.Time    `json:"acceptedAt,omitempty"`
	ArrivedAt     *time.Time    `json:"arrivedAt,omitempty"`
	CompletedAt   *time.Time    `json:"completedAt,omitempty"`
	ETAMinutes    int           `json:"etaMinutes,omitempty"`
	SmsSent       bool          `json:"smsSent"`
	PushSent      bool          `json:"pushSent"`
	DeclineReason string        `json:"reason,omitempty"`
	Metrics       *MatchMetrics `json:"metrics,omitempty"`
}

type MatchMetrics struct {
	DistanceKm float64 `json:"distanceKm"`
	TravelTime int     `json:"travelTime"`
}

type DeclineMatchInput struct {
	Reason string `json:"reason" binding:"required"`
}

type MatchStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
