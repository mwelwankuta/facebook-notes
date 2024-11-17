package summaries

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
)

type SummariesHandler struct {
	useCase           SummariesUseCase
	OpenGraphClientID string
}

func NewSummariesHandler(useCase SummariesUseCase, clientId string) *SummariesHandler {
	return &SummariesHandler{
		useCase:           useCase,
		OpenGraphClientID: clientId,
	}
}

func (h *SummariesHandler) CreateSummaryRequestHandler(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "unauthorized"})
	}

	var dto CreateSummaryRequestDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	if err := utils.Validate(dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	request, err := h.useCase.CreateSummaryRequest(dto, user)
	if err != nil {
		switch err {
		case ErrUnauthorized:
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusCreated, request)
}

func (h *SummariesHandler) GetAllSummariesHandler(c echo.Context) error {
	dto := utils.GetPaginationFromQuery(c)
	summaries, err := h.useCase.GetAllSummaries(dto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, summaries)
}

func (h *SummariesHandler) GetAllRequestsHandler(c echo.Context) error {
	dto := utils.GetPaginationFromQuery(c)
	requests, err := h.useCase.GetAllRequests(dto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *SummariesHandler) GetSummaryByIDHandler(c echo.Context) error {
	id := c.Param("id")
	summary, err := h.useCase.GetSummaryByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}

func (h *SummariesHandler) RateSummaryHandler(c echo.Context) error {
	id := c.Param("id")
	var dto RateSummaryDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	if err := h.useCase.RateSummary(id, dto); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rating updated successfully"})
}

func (h *SummariesHandler) ModerateSummaryHandler(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "unauthorized"})
	}

	id := c.Param("id")
	var dto ModerateRequestDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	if err := utils.Validate(dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	err = h.useCase.ModerateSummary(id, dto, user)
	if err != nil {
		switch err {
		case ErrNotModerator:
			return c.JSON(http.StatusForbidden, utils.ErrorResponse{Error: err.Error()})
		case ErrInvalidStatus:
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Summary moderated successfully"})
}

// EditSummaryHandler handles requests to edit a summary's content
func (h *SummariesHandler) EditSummaryHandler(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "unauthorized"})
	}

	id := c.Param("id")
	var dto EditSummaryDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	if err := utils.Validate(dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	err = h.useCase.EditSummary(id, dto, user)
	if err != nil {
		switch err {
		case ErrNotModerator:
			return c.JSON(http.StatusForbidden, utils.ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Summary edited successfully"})
}

// AddResourceLinkHandler handles requests to add a resource link to a summary
func (h *SummariesHandler) AddResourceLinkHandler(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "unauthorized"})
	}

	summaryID := c.Param("id")
	var dto ResourceLinkDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	if err := utils.Validate(dto); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
	}

	err = h.useCase.AddResourceLink(summaryID, dto, user)
	if err != nil {
		switch err {
		case ErrNotModerator:
			return c.JSON(http.StatusForbidden, utils.ErrorResponse{Error: err.Error()})
		case ErrInvalidResource:
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Resource link added successfully"})
}

// RemoveResourceLinkHandler handles requests to remove a resource link from a summary
func (h *SummariesHandler) RemoveResourceLinkHandler(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "unauthorized"})
	}

	linkID := c.Param("linkId")
	err = h.useCase.RemoveResourceLink(linkID, user)
	if err != nil {
		switch err {
		case ErrNotModerator:
			return c.JSON(http.StatusForbidden, utils.ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Resource link removed successfully"})
}
