package summaries

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
	"gorm.io/gorm"
)

type SummariesRepository struct {
	db *gorm.DB
}

func NewSummariesRepository(db *gorm.DB) *SummariesRepository {
	return &SummariesRepository{db: db}
}

func (r *SummariesRepository) CreateSummary(summary Summary) (Summary, error) {
	summary.ID = uuid.New().String()
	result := r.db.Create(&summary)
	return summary, result.Error
}

func (r *SummariesRepository) CreateSummaryRequest(req SummaryRequest) (SummaryRequest, error) {
	req.ID = uuid.New().String()
	result := r.db.Create(&req)
	return req, result.Error
}

func (r *SummariesRepository) GetAllSummaries(dto models.PaginateDto) ([]Summary, error) {
	var summaries []Summary
	result := r.db.Limit(dto.Limit).Offset(dto.Offset).Find(&summaries)
	return summaries, result.Error
}

func (r *SummariesRepository) GetAllRequests(dto models.PaginateDto) ([]SummaryRequest, error) {
	var requests []SummaryRequest
	result := r.db.Limit(dto.Limit).Offset(dto.Offset).Find(&requests)
	return requests, result.Error
}

func (r *SummariesRepository) GetSummaryByID(id string) (Summary, error) {
	var summary Summary
	result := r.db.First(&summary, "id = ?", id)
	if result.Error != nil {
		return Summary{}, errors.New("summary not found")
	}
	return summary, nil
}

func (r *SummariesRepository) UpdateSummaryRating(id string, rating float64) error {
	result := r.db.Model(&Summary{}).Where("id = ?", id).Update("rating", rating)
	return result.Error
}

func (r *SummariesRepository) UpdateSummaryStatus(id string, status string, moderatorID string, notes string) error {
	now := time.Now()
	return r.db.Model(&Summary{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":          status,
		"moderator_id":    moderatorID,
		"moderated_at":    now,
		"moderator_notes": notes,
	}).Error
}

func (r *SummariesRepository) UpdateAIResponse(id string, aiResponse string) error {
	return r.db.Model(&Summary{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      StatusAIReviewed,
		"ai_response": aiResponse,
	}).Error
}

func (r *SummariesRepository) CreateSummaryEdit(edit SummaryEdit) error {
	return r.db.Create(&edit).Error
}

func (r *SummariesRepository) GetSummaryEdits(summaryID string) ([]SummaryEdit, error) {
	var edits []SummaryEdit
	result := r.db.Where("summary_id = ?", summaryID).Order("version desc").Find(&edits)
	return edits, result.Error
}

func (r *SummariesRepository) AddResourceLink(link ResourceLink) error {
	return r.db.Create(&link).Error
}

func (r *SummariesRepository) RemoveResourceLink(id string) error {
	return r.db.Delete(&ResourceLink{}, "id = ?", id).Error
}

func (r *SummariesRepository) GetSummaryWithResources(id string) (Summary, error) {
	var summary Summary
	result := r.db.Preload("Resources").Preload("EditHistory").First(&summary, "id = ?", id)
	if result.Error != nil {
		return Summary{}, errors.New("summary not found")
	}
	return summary, nil
}

func (r *SummariesRepository) UpdateSummaryContent(id string, content string, version int) error {
	return r.db.Model(&Summary{}).Where("id = ?", id).Updates(map[string]interface{}{
		"content":         content,
		"current_version": version,
	}).Error
}
