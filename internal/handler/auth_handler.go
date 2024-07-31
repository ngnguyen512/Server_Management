package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"server-management/internal/user"
	"server-management/pkg/encryptoha"
	"server-management/pkg/jwtha"
	"server-management/pkg/postgresha"
)

type AuthHandler struct {
	encryptor      encryptoha.IHashEncryptor
	jwtTokenizer   jwtha.JwtTokenizer
	userRepository *postgresha.Repository[user.User]
}

func NewAuthHandler(encryptor encryptoha.IHashEncryptor, jwtTokenizer jwtha.JwtTokenizer, userRepository *postgresha.Repository[user.User]) *AuthHandler {
	return &AuthHandler{
		encryptor:      encryptor,
		jwtTokenizer:   jwtTokenizer,
		userRepository: userRepository,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	hashedPassword, err := h.encryptor.Hash(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing password")
	}

	newUser := user.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	createdUser, err := h.userRepository.CreateOne(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error creating user")
	}

	return c.JSON(http.StatusCreated, createdUser)

}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	foundUser, err := h.userRepository.FindOneByAttribute("username", req.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	match, err := h.encryptor.Compare(foundUser.Password, req.Password)
	if err != nil || !match {
		return c.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	tokenData := map[string]interface{}{"username": foundUser.Username}
	token, err := h.jwtTokenizer.Gencode(tokenData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating token")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
