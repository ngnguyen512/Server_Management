package handler

import (
	"encoding/json"
	"net/http"
	"server-management/internal/user"
	"server-management/pkg/encryptoha"
	"server-management/pkg/jwtha"
	"server-management/pkg/repositories"
)

type AuthHandler struct {
	encryptor      encryptoha.IHashEncryptor
	jwtTokenizer   jwtha.JwtTokenizer
	userRepository repositories.Repository[user.User]
}

func NewAuthHandler(encryptor encryptoha.IHashEncryptor, jwtTokenizer jwtha.JwtTokenizer) *AuthHandler {
	return &AuthHandler{
		encryptor:    encryptor,
		jwtTokenizer: jwtTokenizer,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.encryptor.Hash(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := user.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := h.userRepository.CreateOne(user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := handler.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	match, err := h.encryptor.Compare(user.Password, req.Password)
	if err != nil || !match {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenData := map[string]interface{}{"username": user.Username}
	token, err := h.jwtTokenizer.Gencode(tokenData)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
