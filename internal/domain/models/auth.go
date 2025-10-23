package models

import "time"

type OTPRequest struct {
	Phone   string `json:"phone" binding:"required"`
	Locale  string `json:"locale,omitempty" binding:"omitempty"`
	Channel string `json:"channel,omitempty" binding:"omitempty,oneof=sms call"`
}

type OTPRequestResponse struct {
	ExpiresIn         int `json:"expiresIn"`
	AttemptsRemaining int `json:"attemptsRemaining"`
}

type OTPVerification struct {
	Phone    string `json:"phone" binding:"required"`
	OTP      string `json:"otp" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	DeviceID string `json:"deviceId" binding:"required"`
}

type Session struct {
	Token        string    `json:"sessionToken"`
	RefreshToken string    `json:"refreshToken"`
	User         User      `json:"user"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
