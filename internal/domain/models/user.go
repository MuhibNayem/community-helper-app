package models

import "time"

type User struct {
	ID                string            `json:"id"`
	Phone             string            `json:"phone"`
	Name              string            `json:"name"`
	PhotoURL          string            `json:"photoUrl,omitempty"`
	Language          string            `json:"language"`
	IsHelper          bool              `json:"isHelper"`
	HelperStatus      string            `json:"helperStatus,omitempty"`
	NotificationPrefs NotificationPrefs `json:"notificationPrefs"`
	KYCStatus         string            `json:"kycStatus,omitempty"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

type NotificationPrefs struct {
	QuietHours QuietHours `json:"quietHours"`
	UrgentSMS  bool       `json:"urgentSms"`
}

type QuietHours struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type HelperProfile struct {
	UserID       string       `json:"userId"`
	Skills       []string     `json:"skills"`
	Availability Availability `json:"availability"`
	Rating       float64      `json:"rating"`
	RatingCount  int          `json:"ratingCount"`
	Badges       []string     `json:"badges"`
	OptedIn      bool         `json:"optedIn"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type Availability struct {
	Weekly     []AvailabilitySlot      `json:"weekly"`
	Exceptions []AvailabilityException `json:"exceptions,omitempty"`
}

type AvailabilitySlot struct {
	Day   string `json:"day"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityException struct {
	Date  string          `json:"date"`
	Slots []ExceptionSlot `json:"slots"`
}

type ExceptionSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type KYCDocument struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	DocumentType string    `json:"documentType"`
	FileURL      string    `json:"fileUrl"`
	Status       string    `json:"status"`
	SubmittedAt  time.Time `json:"submittedAt"`
	ReviewedAt   time.Time `json:"reviewedAt,omitempty"`
	ReviewerID   string    `json:"reviewerId,omitempty"`
	Notes        string    `json:"notes,omitempty"`
}

type UserUpdate struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=1"`
	Language *string `json:"language,omitempty" binding:"omitempty,len=2"`
	PhotoURL *string `json:"photoUrl,omitempty" binding:"omitempty,url"`
}

type HelperToggle struct {
	OptedIn bool `json:"optedIn"`
}

type SkillsUpdate struct {
	Skills []string `json:"skills" binding:"required"`
}

type AvailabilityUpdate struct {
	Weekly     []AvailabilitySlotInput      `json:"weekly" binding:"required"`
	Exceptions []AvailabilityExceptionInput `json:"exceptions,omitempty"`
}

type AvailabilitySlotInput struct {
	Day   string `json:"day" binding:"required"`
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

type AvailabilityExceptionInput struct {
	Date  string                  `json:"date" binding:"required"`
	Slots []AvailabilitySlotInput `json:"slots" binding:"required"`
}

type KYCDocumentUpload struct {
	DocumentType string `json:"documentType" binding:"required"`
	FileURL      string `json:"fileUrl" binding:"required,url"`
}
