package summaries

type SummariesHandler struct {
	useCase           SummariesUseCase
	OpenGraphClientID string
}

// NewSummariesHandler creates a new AuthHandler
func NewSummariesHandler(useCase SummariesUseCase, clientId string) *SummariesHandler {
	return &SummariesHandler{
		useCase:           useCase,
		OpenGraphClientID: clientId,
	}
}
