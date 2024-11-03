package summaries

import (
	"github.com/mwelwankuta/facebook-notes/pkg/config"
)

type SummariesUseCase struct {
	config config.Config
	repo   SummariesRepository
}

func NewSummariesUseCase(repo SummariesRepository, cfg config.Config) *SummariesUseCase {
	return &SummariesUseCase{
		repo:   repo,
		config: cfg,
	}
}
