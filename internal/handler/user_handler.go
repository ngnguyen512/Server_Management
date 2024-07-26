package handler

import (
	"net/http"
	"server-management/internal/user" // Adjust according to your project structure
	"server-management/pkg/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4" // Assuming Echo V4 is used for handling HTTP requests
)

type UserHandler struct {
	UserRepository repositories.Repository[user.User]
}

func NewUserHandler(userRepo repositories.Repository[user.User]) *UserHandler {
	return &UserHandler{
		UserRepository: userRepo,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user user.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid User Data")
	}

	newUser, err := h.UserRepository.CreateOne(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create user")
	}
	return c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) GetUserById(c echo.Context) error {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	user, err := h.UserRepository.FindOneById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid update data")
	}
	updatedUser, err := h.UserRepository.UpdateOneById(id, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user")
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid UUID format")
	}
	if err := h.UserRepository.DeleteOneById(id); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.NoContent(http.StatusNoContent)
}
