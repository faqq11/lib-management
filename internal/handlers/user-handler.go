package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/faqq11/lib-management/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	DB *sqlx.DB
}

func (userHandler *UserHandler) Register(writer http.ResponseWriter, request *http.Request) {
	var userReq struct {
		Username string
		Password string
	}

	err := json.NewDecoder(request.Body).Decode(&userReq)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
			return
	}

	if userReq.Username == "" || userReq.Password == "" {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Username and password are required")
			return
	}

	hashed, err := helper.HashPassword(userReq.Password)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Hash error")
		return
	}

	_, err = userHandler.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", userReq.Username, hashed)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Username already exist")
		return
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}

func (userHandler *UserHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var userReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&userReq)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
			return
	}

	if userReq.Username == "" || userReq.Password == "" {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Username and password are required")
			return
	}

	var user models.User
	err = userHandler.DB.Get(&user, "SELECT id, username, password, role FROM users WHERE username=$1", userReq.Username)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "invalid credentials")
			return
	}

	errPasswordCheck := helper.CheckPassword(user.Password, userReq.Password)
	if errPasswordCheck != nil {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "invalid credentials")
			return
	}

	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
			return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
			"access_token": token,
	})
}