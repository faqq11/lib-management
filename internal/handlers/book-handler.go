package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		log.Printf("InsertBook - JSON decode error: %v", err)
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
			log.Printf("InsertBook - Insert error: %v", err)
			helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
			"message": "Book created successfully",
		})
		return

	} else if err != nil {
		log.Printf("InsertBook - Database error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = bookHandler.DB.Exec(`
			UPDATE books
			SET stock = stock + 1
			WHERE title = $1
    `, bookInput.Title)
	if err != nil {
		log.Printf("InsertBook - Update stock error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to increase stock")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"message": "Book already exists. Stock increased by 1",
	})
}

func (bookHandler *BookHandler) GetAllBooks(writer http.ResponseWriter, request *http.Request) {
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
		log.Printf("GetAllBooks - Select error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, books)
}

func (bookHandler *BookHandler) GetBookById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("GetBookById - Invalid ID: %s, error: %v", id, err)
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

		log.Printf("GetBookById - Database error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, book)
}

func (bookHandler *BookHandler) UpdateBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("UpdateBook - Invalid ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var book models.Book

	err = json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		log.Printf("UpdateBook - JSON decode error: %v", err)
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
		log.Printf("UpdateBook - Update error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update book")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("UpdateBook - RowsAffected error: %v", err)
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

func (bookHandler *BookHandler) IncreaseStock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("IncreaseStock - Invalid ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	result, err := bookHandler.DB.Exec(`
		UPDATE books
		SET stock = stock + 1
		WHERE id = $1
	`, bookId)

	if err != nil {
		log.Printf("IncreaseStock - Update error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update stock")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("IncreaseStock - RowsAffected error: %v", err)
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

func (bookHandler *BookHandler) DecreaseStock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("DecreaseStock - Invalid ID: %s, error: %v", id, err)
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
		log.Printf("DecreaseStock - Select stock error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch stock")
		return
	}

	if stock <= 0 {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Stock is already 0, cannot decrease")
		return
	}

	result, err := bookHandler.DB.Exec(`
		UPDATE books
    SET stock = stock - 1
    WHERE id = $1
    `, bookId)

	if err != nil {
		log.Printf("DecreaseStock - Update error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update stock")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DecreaseStock - RowsAffected error: %v", err)
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

func (bookHandler *BookHandler) DeleteBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("DeleteBook - Invalid ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	result, err := bookHandler.DB.Exec(`DELETE FROM books WHERE id = $1`, bookId)
	if err != nil {
		log.Printf("DeleteBook - Delete error: %v", err)
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

func (bookHandler *BookHandler) SearchBooks(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Query().Get("title")
	categoryID := request.URL.Query().Get("category_id")

	var books []response.BookResponse
	var args []interface{}
	var conditions []string

	baseQuery := `
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
	`

	argIndex := 1

	if title != "" {
		conditions = append(conditions, "b.title ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+title+"%")
		argIndex++
	}

	if categoryID != "" {
		categoryIDInt, err := strconv.Atoi(categoryID)
		if err != nil {
			log.Printf("SearchBooks - Invalid category ID: %s, error: %v", categoryID, err)
			helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid category ID")
			return
		}
		conditions = append(conditions, "b.category_id = $"+strconv.Itoa(argIndex))
		args = append(args, categoryIDInt)
		argIndex++
	}

	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	finalQuery += " ORDER BY b.title"

	err := bookHandler.DB.Select(&books, finalQuery, args...)
	if err != nil {
		log.Printf("SearchBooks - Select error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if len(books) == 0 {
		helper.SuccessResponse(writer, http.StatusOK, map[string]string{
			"message": "No books found matching your search criteria",
		})
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, books)
}
