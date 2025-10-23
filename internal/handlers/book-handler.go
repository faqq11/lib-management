package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (bookHandler *BookHandler) InsertBook(writer http.ResponseWriter, request *http.Request) {
	var bookInput models.Book

	err := json.NewDecoder(request.Body).Decode(&bookInput)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if bookInput.Title == "" {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Title is required")
		return
	}

	var title string
	err = bookHandler.DB.Get(&title, `SELECT title FROM books WHERE title = $1`, bookInput.Title)

	if errors.Is(err, sql.ErrNoRows) {
		_, err = bookHandler.DB.Exec(`
			INSERT INTO books (title, author, category_id, stock)
			VALUES ($1, $2, $3, $4)
		`, bookInput.Title, bookInput.Author, bookInput.CategoryID, bookInput.Stock)
		if err != nil {
			helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
			"message": "Book created successfully",
		})
		return

	} else if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = bookHandler.DB.Exec(`
		UPDATE books
		SET stock = stock + 1
		WHERE title = $1
	`, bookInput.Title)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to increase stock")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"message": "Book already exists. Stock increased by 1",
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
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
			return
		}

		helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, book)
}

func (bookHandler *BookHandler) UpdateBook(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return 
	}

	var book models.Book

	err = json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	result, err := bookHandler.DB.Exec(`
	UPDATE books
	SET title = $1,
		author = $2,
		category_id = $3,
		stock = $4
	WHERE id = $5
	`, book.Title, book.Author, book.CategoryID, book.Stock, bookId)

	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update book")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check update result")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]string{
		"message": "Book updated successfully",
	})
}

func (bookHandler *BookHandler) IncreaseStock(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	result, err := bookHandler.DB.Exec(`
		UPDATE books
		SET stock = stock + 1
		WHERE id = $1
	`, bookId)

	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update stock")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check update result")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]string{
		"message": "Stock increased by 1",
	})
}

func (bookHandler *BookHandler) DecreaseStock(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var stock int

	err = bookHandler.DB.Get(&stock, `SELECT stock FROM books WHERE id = $1`, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
			return
		}
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch stock")
		return
	}

	if stock <= 0{
		helper.ErrorResponse(writer, http.StatusBadRequest, "Stock is already 0, cannot decrease")
		return 
	}

	result, err := bookHandler.DB.Exec(`
		UPDATE books
    SET stock = stock - 1
    WHERE id = $1
	`, bookId)

	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update stock")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check update result")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]string{
		"message": "Stock decreased by 1",
	})
}

func (bookHandler *BookHandler) DeleteBook(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	result, err:=bookHandler.DB.Exec(`DELETE FROM books WHERE id = $1`, bookId)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check deleted book")
		return
	}

  helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"message": "Book deleted successfully",
	})
}