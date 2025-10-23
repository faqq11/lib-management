package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/faqq11/lib-management/internal/helper"
	"github.com/faqq11/lib-management/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type BorrowHandler struct {
	DB *sqlx.DB
}

func (borrowHandler *BorrowHandler) BorrowBook(writer http.ResponseWriter, request *http.Request){
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
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid book ID")
		return 
	}

	var stock int
	err = borrowHandler.DB.Get(&stock, `SELECT stock FROM books where id = $1`, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Book not found")
			return
		}
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch stock")
		return
	}

	if stock <= 0{
		helper.ErrorResponse(writer, http.StatusBadRequest, "Book is not available")
		return 
	}

	_, err = borrowHandler.DB.Exec(`INSERT INTO borrowings (user_id, book_id) VALUES ($1, $2)`, userId, bookId)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = borrowHandler.DB.Exec(`UPDATE books SET stock = stock - 1 WHERE id = $1`, bookId)
	if err != nil{
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	helper.SuccessResponse(writer, http.StatusCreated, map[string]string{
		"message": "Book borrowed successfully",
	})
}

func (borrowHandler *BorrowHandler) ReturnBook(writer http.ResponseWriter, request *http.Request){
	vars := mux.Vars(request)
	id := vars["id"]

	borrowId, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusBadRequest, "Invalid borrow ID")
		return 
	}

	var bookId int
	err = borrowHandler.DB.Get(&bookId, `SELECT book_id FROM borrowings WHERE id = $1`, borrowId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.ErrorResponse(writer, http.StatusNotFound, "Borrowing data not found")
			return
		}
		helper.ErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch book ID")
		return
	}

	_, err = borrowHandler.DB.Exec(`UPDATE borrowings SET returned_at = $1 WHERE id = $2`, time.Now() , borrowId)
	if err != nil {
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = borrowHandler.DB.Exec(`UPDATE books SET stock = stock + 1 WHERE id = $1`, bookId)
	if err != nil{
		helper.ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(writer, http.StatusOK, map[string]string{
		"message": "Book returned successfully",
	})
}