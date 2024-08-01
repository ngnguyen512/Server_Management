package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"server-management/internal/user"
	"server-management/pkg/encryptoha"
	"server-management/pkg/jwtha"
	"server-management/pkg/repositories"
	"server-management/pkg/validatorha" // Adjust the import path as needed
	"time"
)

// Initialize the validator
var validator = validatorha.NewValidator(validatorha.ValidatorConfig{
	EnableIPv4Validation: true,
	CustomMessages: map[string]string{
		"required": "The {0} field is required.",
		"uuid":     "The {0} field must be a valid UUID.",
		"ipv4":     "The {0} field must be a valid IPv4 address.",
	},
})

type AuthHandler struct {
	encryptor      encryptoha.IHashEncryptor
	jwtTokenizer   jwtha.JwtTokenizer
	userRepository repositories.Repository[user.User]
}

func NewAuthHandler(encryptor encryptoha.IHashEncryptor, jwtTokenizer jwtha.JwtTokenizer, userRepository repositories.Repository[user.User]) *AuthHandler {
	return &AuthHandler{
		encryptor:      encryptor,
		jwtTokenizer:   jwtTokenizer,
		userRepository: userRepository,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// Validate the request
	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"validation_error": err.Error()})
	}

	hashedPassword, err := h.encryptor.Hash(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing password")
	}

	newUser := user.User{
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := h.userRepository.CreateOne(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error creating user")
	}

	return c.JSON(http.StatusCreated, createdUser)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// Validate the request
	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"validation_error": err.Error()})
	}

	foundUser, err := h.userRepository.FindOneByAttribute("username", req.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	match, err := h.encryptor.Compare(foundUser.Password, req.Password)
	if err != nil || !match {
		return c.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	// Generate token only if the password matches
	tokenData := map[string]interface{}{"username": foundUser.Username}
	token, err := h.jwtTokenizer.Gencode(tokenData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating token")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
