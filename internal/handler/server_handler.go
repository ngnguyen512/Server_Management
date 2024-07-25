package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-management/internal/server"
	"server-management/pkg/repositories"
)

type ServerHandler struct {
	ServerRepository repositories.Repository[server.Server]
}

func NewServerHandler(repo repositories.Repository[server.Server]) *ServerHandler {
	return &ServerHandler{ServerRepository: repo}
}

func (h *ServerHandler) GetServerById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid server ID")
	}
	server, err := h.ServerRepository.FindOneById(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Server not found")
	}
	return c.JSON(http.StatusOK, server)
}
func (h *ServerHandler) CreateServer(c echo.Context) error {
	var server server.Server
	if err := c.Bind(&server); err != nil {
		return c.String(http.StatusBadRequest, "Error reading server data")
	}
	createdServer, err := h.ServerRepository.CreateOne(server)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating server")
	}

	return c.JSON(http.StatusCreated, createdServer)
}

func (h *ServerHandler) UpdateServer(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid server ID")
	}
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.String(http.StatusBadRequest, "Error reading update data")
	}
	updatedServer, err := h.ServerRepository.UpdateOneById(id, updates)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating server")
	}
	return c.JSON(http.StatusOK, updatedServer)
}

func (h *ServerHandler) DeleteOneById(c echo.Context) error {
	// Retrieve the server ID as a string from the route parameter
	serverID := c.Param("id")

	// Convert the string ID to a uuid.UUID type
	id, err := uuid.Parse(serverID)
	if err != nil {
		// If there is an error in parsing, it means the UUID format is not valid
		return c.JSON(http.StatusBadRequest, "Invalid UUID format")
	}

	// Now pass the correctly typed UUID to the repository method
	err = h.ServerRepository.DeleteOneById(id)
	if err != nil {
		// If deletion fails, return an internal server error
		return c.JSON(http.StatusInternalServerError, "Failed to delete server")
	}

	// If no error, return no content to indicate successful deletion
	return c.NoContent(http.StatusNoContent)
}
