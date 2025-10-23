package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/jmoiron/sqlx"
)

type BookHandler struct {
	DB *sqlx.DB
}

func (bookHandler *BookHandler) InsertBook(writer http.ResponseWriter, request *http.Request){
	var bookInput struct{
		Title string
		Author string
		CategoryId *int
		Stock *int
	}

	err := json.NewDecoder(request.Body).Decode(&bookInput)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
			return
	}

	if bookInput.Title == ""{
			helper.ErrorResponse(writer, http.StatusBadRequest, "Title is required")
			return
	}

	_, err = bookHandler.DB.Exec("INSERT INTO books (title, author, category_id, stock) VALUES ($1, $2, $3, $4)", bookInput.Title, bookInput.Author, bookInput.CategoryId, bookInput.Stock)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
		"message": "Books created successfully",
	})
}