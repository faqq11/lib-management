package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type CategoryHandler struct {
	DB *sqlx.DB
}

func (categoryHandler *CategoryHandler) CreateCategory(writer http.ResponseWriter, request *http.Request) {
	var categoryInput struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(request.Body).Decode(&categoryInput)
	if err != nil {
		log.Printf("CreateCategory - JSON decode error: %v", err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if categoryInput.Name == "" {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Category name required")
		return
	}

	_, err = categoryHandler.DB.Exec("INSERT INTO categories (name) VALUES ($1)", categoryInput.Name)
	if err != nil {
		log.Printf("CreateCategory - Insert error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
		"message": "Category created successfully",
	})
}

func (categoryHandler *CategoryHandler) DeleteCategory(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	categoryId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("DeleteCategory - Invalid category ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid category ID")
		return
	}

	result, err := categoryHandler.DB.Exec("DELETE FROM categories WHERE id = $1", categoryId)
	if err != nil {
		log.Printf("DeleteCategory - Delete error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Category not found")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"message": "Category deleted successfully",
	})
}
