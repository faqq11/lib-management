package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/faqq11/lib-management/internal/helper"
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

	json.NewDecoder(request.Body).Decode(&userReq)

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