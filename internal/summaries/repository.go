package summaries

import (
	"gorm.io/gorm"
)

type SummariesRepository struct {
	db *gorm.DB
}

func NewSummariesRepository(db *gorm.DB) *SummariesRepository {
	return &SummariesRepository{
		db: db,
	}
}
