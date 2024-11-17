package summaries

import (
	"time"

	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

const (
	StatusPending    = "pending"
	StatusAIReviewed = "ai_reviewed"
	StatusApproved   = "approved"
	StatusRejected   = "rejected"

	RoleModerator = "moderator"
	RoleUser      = "user"
)

type Summary struct {
	ID             string    `json:"id" gorm:"primarykey"`
	Content        string    `json:"content"`
	Summary        string    `json:"summary"`
	IsVerified     bool      `json:"is_verified"`
	IsSummarizedAI bool      `json:"is_summarized_ai"`
	Rating         float64   `json:"rating"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         string    `json:"user_id"`
	User           models.User
	Status         string         `json:"status"`
	ModeratorID    *string        `json:"moderator_id,omitempty"`
	ModeratedAt    *time.Time     `json:"moderated_at,omitempty"`
	AIResponse     string         `json:"ai_response,omitempty"`
	ModeratorNotes string         `json:"moderator_notes,omitempty"`
	Resources      []ResourceLink `json:"resources"`
	EditHistory    []SummaryEdit  `json:"edit_history"`
	CurrentVersion int            `json:"current_version" gorm:"default:1"`
}

type SummaryRequest struct {
	ID        string    `json:"id" gorm:"primarykey"`
	Content   string    `json:"content"`
	Metadata  string    `json:"metadata"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	User      models.User
}

type CreateSummaryRequestDto struct {
	Content  string `json:"content" validate:"required,min=10"`
	Metadata string `json:"metadata" validate:"required"`
}

type RateSummaryDto struct {
	Rating float64 `json:"rating" validate:"required,min=0,max=5"`
}

type ModerateRequestDto struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Notes  string `json:"notes" validate:"required_if=Action reject"`
}

type ResourceLink struct {
	ID          string    `json:"id" gorm:"primarykey"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SummaryID   string    `json:"summary_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

type SummaryEdit struct {
	ID          string    `json:"id" gorm:"primarykey"`
	SummaryID   string    `json:"summary_id"`
	Content     string    `json:"content"`
	EditedBy    string    `json:"edited_by"`
	EditedAt    time.Time `json:"edited_at"`
	Version     int       `json:"version"`
	EditMessage string    `json:"edit_message"`
}

type EditSummaryDto struct {
	Content     string `json:"content" validate:"required,min=10"`
	EditMessage string `json:"edit_message" validate:"required,min=5"`
}

type ResourceLinkDto struct {
	URL         string `json:"url" validate:"required,url"`
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty,min=10"`
}
