package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/faqq11/lib-management/internal/middleware"
	"github.com/faqq11/lib-management/internal/models/response"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type BorrowHandler struct {
	DB *sqlx.DB
}

func (borrowHandler *BorrowHandler) BorrowBook(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(middleware.UserContextKey)
	if user == nil {
		helper.ErrorResponse(writer, http.StatusUnauthorized, "User context not found")
		return
	}

	userClaims := user.(middleware.UserClaims)
	userId := userClaims.UserID

	vars := mux.Vars(request)
	id := vars["id"]

	bookId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("BorrowBook - Invalid book ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return
	}

	tx, err := borrowHandler.DB.Beginx()
	if err != nil {
		log.Printf("BorrowBook - Transaction start error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var existingBorrow int
	err = tx.Get(&existingBorrow, `
			SELECT COUNT(*) FROM borrowings 
			WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
    `, userId, bookId)
	if err != nil {
		log.Printf("BorrowBook - Check existing borrow error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check existing borrow")
		return
	}

	if existingBorrow > 0 {
		helper.ErrorResponse(writer, http.StatusBadRequest, "You have already borrowed this book")
		return
	}

	var stock int
	err = tx.Get(&stock, `SELECT stock FROM books WHERE id = $1`, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
			return
		}
		log.Printf("BorrowBook - Fetch stock error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch stock")
		return
	}

	if stock <= 0 {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Book is not available")
		return
	}

	_, err = tx.Exec(`
			INSERT INTO borrowings (user_id, book_id, borrowed_at) 
			VALUES ($1, $2, $3)
    `, userId, bookId, time.Now())
	if err != nil {
		log.Printf("BorrowBook - Insert borrowing error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to create borrowing record")
		return
	}

	result, err := tx.Exec(`UPDATE books SET stock = stock - 1 WHERE id = $1`, bookId)
	if err != nil {
		log.Printf("BorrowBook - Update stock error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update book stock")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("BorrowBook - RowsAffected error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check update result")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("BorrowBook - Transaction commit error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]interface{}{
		"message": "Book borrowed successfully",
	})
}

func (borrowHandler *BorrowHandler) ReturnBook(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(middleware.UserContextKey)
	if user == nil {
		helper.ErrorResponse(writer, http.StatusUnauthorized, "User context not found")
		return
	}

	userClaims := user.(middleware.UserClaims)
	userId := userClaims.UserID

	vars := mux.Vars(request)
	id := vars["id"]

	borrowId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("ReturnBook - Invalid borrow ID: %s, error: %v", id, err)
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid borrow ID")
		return
	}

	tx, err := borrowHandler.DB.Beginx()
	if err != nil {
		log.Printf("ReturnBook - Transaction start error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var borrowData struct {
		ID     int `db:"id"`
		BookID int `db:"book_id"`
		UserID int `db:"user_id"`
	}

	err = tx.Get(&borrowData, `
		SELECT id, book_id, user_id 
		FROM borrowings 
		WHERE id = $1 AND returned_at IS NULL
	`, borrowId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Borrowing record not found or already returned")
			return
		}
		log.Printf("ReturnBook - Fetch borrowing data error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch borrowing data")
		return
	}

	if borrowData.UserID != userId {
		helper.ErrorResponse(writer, http.StatusForbidden, "You can only return your own borrowed books")
		return
	}

	result, err := tx.Exec(`
		UPDATE borrowings 
		SET returned_at = $1 
		WHERE id = $2
	`, time.Now(), borrowId)
	if err != nil {
		log.Printf("ReturnBook - Update borrowing error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update borrowing record")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ReturnBook - RowsAffected error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to check update result")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorResponse(writer, http.StatusNotFound, "Borrowing record not found")
		return
	}

	_, err = tx.Exec(`UPDATE books SET stock = stock + 1 WHERE id = $1`, borrowData.BookID)
	if err != nil {
		log.Printf("ReturnBook - Update stock error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to update book stock")
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("ReturnBook - Transaction commit error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]interface{}{
		"message": "Book returned successfully",
	})
}

func (borrowHandler *BorrowHandler) GetUserBorrowings(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(middleware.UserContextKey)
	if user == nil {
		helper.ErrorResponse(writer, http.StatusUnauthorized, "User context not found")
		return
	}

	userClaims := user.(middleware.UserClaims)
	userId := userClaims.UserID

	var borrowings []response.UserBorrowingResponse
	err := borrowHandler.DB.Select(&borrowings, `
		SELECT 
			br.id,
			br.book_id,
			b.title as book_title,
			b.author,
			br.borrowed_at,
			br.returned_at,
		CASE 
			WHEN br.returned_at IS NULL THEN 'borrowed'
			ELSE 'returned'
		END as status
		FROM borrowings br
		JOIN books b ON br.book_id = b.id
		WHERE br.user_id = $1
		ORDER BY br.borrowed_at DESC
	`, userId)

	if err != nil {
		log.Printf("GetUserBorrowings - Select error: %v", err)
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch borrowings")
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, borrowings)
}
