package handler

import (
	"net/http"
	"server-management/internal/health_event" // Adjust according to your project structure
	"server-management/pkg/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4" // Assuming Echo V4 is used for handling HTTP requests
)

type HealthHandler struct {
	HealthRepository repositories.Repository[health_event.HealthEvent]
}

func NewHealthHandler(healthEventRepo repositories.Repository[health_event.HealthEvent]) *HealthHandler {
	return &HealthHandler{
		HealthRepository: healthEventRepo,
	}
}

func (h *HealthHandler) CreateHealthEvent(c echo.Context) error {
	var event health_event.HealthEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Health Event Data")
	}

	newEvent, err := h.HealthRepository.CreateOne(event)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create health event")
	}
	return c.JSON(http.StatusCreated, newEvent)
}

func (h *HealthHandler) GetHealthEventById(c echo.Context) error {
	eventID := c.Param("id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	event, err := h.HealthRepository.FindOneById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Health event not found")
	}
	return c.JSON(http.StatusOK, event)
}

func (h *HealthHandler) UpdateHealthEvent(c echo.Context) error {
	eventID := c.Param("id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid update data")
	}
	updatedEvent, err := h.HealthRepository.UpdateOneById(id, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update health event")
	}
	return c.JSON(http.StatusOK, updatedEvent)
}

func (h *HealthHandler) DeleteHealthEvent(c echo.Context) error {
	eventID := c.Param("id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	if err := h.HealthRepository.DeleteOneById(id); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete health event")
	}
	return c.NoContent(http.StatusNoContent)
}
