package handlers

import (
	"encoding/json"
	"log"
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
		Username string  `json:"username"`
		Password string  `json:"password"`
		Role     *string `json:"role,omitempty"`
	}

	err := json.NewDecoder(request.Body).Decode(&userReq)
	if err != nil {
		log.Printf("Register - JSON decode error: %v", err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if userReq.Username == "" || userReq.Password == "" {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Username and password are required")
		return
	}

	hashed, err := helper.HashPassword(userReq.Password)
	if err != nil {
		log.Printf("Register - Hash password error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	if userReq.Role != nil {
		_, err = userHandler.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", userReq.Username, hashed, userReq.Role)
		if err != nil {
			log.Printf("Register - Insert user with role error: %v", err)
			helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}
		helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
			"message": "User created successfully",
		})
		return
	}

	_, err = userHandler.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", userReq.Username, hashed)
	if err != nil {
		log.Printf("Register - Insert user error: %v", err)
		helper.ErrorResponse(writer, http.StatusConflict, "Username already exist")
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
		log.Printf("Login - JSON decode error: %v", err)
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
		log.Printf("Login - User lookup error: %v", err)
		helper.ErrorResponse(writer, http.StatusUnauthorized, "invalid credentials")
		return
	}

	errPasswordCheck := helper.CheckPassword(user.Password, userReq.Password)
	if errPasswordCheck != nil {
		log.Printf("Login - Password check error: %v", errPasswordCheck)
		helper.ErrorResponse(writer, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("Login - JWT generation error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"access_token": token,
	})
}
