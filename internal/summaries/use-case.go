package summaries

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mwelwankuta/facebook-notes/pkg/adapters"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotModerator    = errors.New("user is not a moderator")
	ErrInvalidStatus   = errors.New("invalid summary status")
	ErrInvalidResource = errors.New("invalid resource link")
)

type SummariesUseCase struct {
	config config.Config
	repo   SummariesRepository
	redis  *adapters.RedisClient
}

func NewSummariesUseCase(repo SummariesRepository, cfg config.Config, redis *adapters.RedisClient) *SummariesUseCase {
	return &SummariesUseCase{
		repo:   repo,
		config: cfg,
		redis:  redis,
	}
}

func (uc *SummariesUseCase) CreateSummaryRequest(dto CreateSummaryRequestDto, user models.User) (SummaryRequest, error) {
	if user.ID == "" {
		return SummaryRequest{}, ErrUnauthorized
	}

	request := SummaryRequest{
		Content:  dto.Content,
		Metadata: dto.Metadata,
		UserID:   user.ID,
		Status:   StatusPending,
	}

	// Create the request
	newRequest, err := uc.repo.CreateSummaryRequest(request)
	if err != nil {
		return SummaryRequest{}, err
	}

	// TODO: Trigger AI summarization asynchronously
	go uc.processAISummarization(newRequest.ID)

	return newRequest, nil
}

func (uc *SummariesUseCase) ModerateSummary(id string, dto ModerateRequestDto, user models.User) error {
	if user.Role != RoleModerator {
		return ErrNotModerator
	}

	var status string
	switch dto.Action {
	case "approve":
		status = StatusApproved
	case "reject":
		status = StatusRejected
	default:
		return ErrInvalidStatus
	}

	return uc.repo.UpdateSummaryStatus(id, status, user.ID, dto.Notes)
}

func (uc *SummariesUseCase) processAISummarization(requestID string) {
	// TODO: Implement AI summarization logic using OpenAI
	// 1. Get the request content
	// 2. Call OpenAI API
	// 3. Update the summary with AI response
	// uc.repo.UpdateAIResponse(requestID, aiResponse)
}

func (uc *SummariesUseCase) GetAllSummaries(dto models.PaginateDto) ([]Summary, error) {
	return uc.repo.GetAllSummaries(dto)
}

func (uc *SummariesUseCase) GetAllRequests(dto models.PaginateDto) ([]SummaryRequest, error) {
	return uc.repo.GetAllRequests(dto)
}

func (uc *SummariesUseCase) GetSummaryByID(id string) (Summary, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("summary:%s", id)

	// Try to get from cache first
	var summary Summary
	err := uc.redis.Get(ctx, cacheKey, &summary)
	if err == nil {
		return summary, nil
	}

	// If not in cache, get from database
	summary, err = uc.repo.GetSummaryByID(id)
	if err != nil {
		return Summary{}, err
	}

	// Cache the result
	uc.redis.Set(ctx, cacheKey, summary, 30*time.Minute)
	return summary, nil
}

func (uc *SummariesUseCase) RateSummary(id string, dto RateSummaryDto) error {
	return uc.repo.UpdateSummaryRating(id, dto.Rating)
}

// EditSummary allows moderators to edit a summary's content and keeps track of edit history
func (uc *SummariesUseCase) EditSummary(id string, dto EditSummaryDto, user models.User) error {
	if user.Role != RoleModerator {
		return ErrNotModerator
	}

	summary, err := uc.repo.GetSummaryWithResources(id)
	if err != nil {
		return err
	}

	// Create edit history entry
	edit := SummaryEdit{
		ID:          uuid.New().String(),
		SummaryID:   summary.ID,
		Content:     dto.Content,
		EditedBy:    user.ID,
		EditedAt:    time.Now(),
		Version:     summary.CurrentVersion + 1,
		EditMessage: dto.EditMessage,
	}

	if err := uc.repo.CreateSummaryEdit(edit); err != nil {
		return err
	}

	// Update summary content
	err = uc.repo.UpdateSummaryContent(id, dto.Content, edit.Version)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("summary:%s", id)
	uc.redis.Delete(ctx, cacheKey)

	return nil
}

// AddResourceLink adds a resource link to a summary
func (uc *SummariesUseCase) AddResourceLink(summaryID string, dto ResourceLinkDto, user models.User) error {
	if user.Role != RoleModerator {
		return ErrNotModerator
	}

	link := ResourceLink{
		ID:          uuid.New().String(),
		URL:         dto.URL,
		Title:       dto.Title,
		Description: dto.Description,
		SummaryID:   summaryID,
		CreatedAt:   time.Now(),
		CreatedBy:   user.ID,
	}

	return uc.repo.AddResourceLink(link)
}

// RemoveResourceLink removes a resource link from a summary
func (uc *SummariesUseCase) RemoveResourceLink(linkID string, user models.User) error {
	if user.Role != RoleModerator {
		return ErrNotModerator
	}
	return uc.repo.RemoveResourceLink(linkID)
}

// GetSummaryWithResources returns a summary with its resources and edit history
func (uc *SummariesUseCase) GetSummaryWithResources(id string) (Summary, error) {
	return uc.repo.GetSummaryWithResources(id)
}
