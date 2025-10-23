package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/faqq11/lib-management/internal/models"
	"github.com/faqq11/lib-management/internal/models/response"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type BookHandler struct {
	DB *sqlx.DB
}

func (bookHandler *BookHandler) InsertBook(writer http.ResponseWriter, request *http.Request){
	var bookInput models.Book

	err := json.NewDecoder(request.Body).Decode(&bookInput)
	if err != nil {
			helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
			return
	}

	if bookInput.Title == ""{
			helper.ErrorResponse(writer, http.StatusBadRequest, "Title is required")
			return
	}

	_, err = bookHandler.DB.Exec("INSERT INTO books (title, author, category_id, stock) VALUES ($1, $2, $3, $4)", bookInput.Title, bookInput.Author, bookInput.CategoryID, bookInput.Stock)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
		"message": "Books created successfully",
	})
}

func (bookHandler *BookHandler) GetAllBooks(writer http.ResponseWriter, request *http.Request){
	var books []response.BookResponse

	err := bookHandler.DB.Select(&books, `
		SELECT 
			b.id, 
			b.title, 
			b.author, 
			b.category_id,
			c.name AS category,
			b.stock, 
			b.created_at
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		ORDER BY b.id;
    `)
    if err != nil {
        helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
        return
    }

    helper.SuccessResponse(writer, http.StatusOK, books)
}

func (bookHandler *BookHandler) GetBookById(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return 
	}

	var book response.BookResponse

	err = bookHandler.DB.Get(&book, `
	SELECT 
		b.id, 
		b.title, 
		b.author, 
		b.category_id,
		c.name AS category,
		b.stock, 
		b.created_at
	FROM books b
	JOIN categories c ON b.category_id = c.id
	WHERE b.id = $1
	`, bookId)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, book)
}