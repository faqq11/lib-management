package models

import "time"

type Borrowing struct {
    ID int `db:"id" json:"id"`
    UserID int `db:"user_id" json:"user_id"`
    BookID int `db:"book_id" json:"book_id"`
    BorrowedAt time.Time `db:"borrowed_at" json:"borrowed_at"`
    ReturnedAt *time.Time `db:"returned_at" json:"returned_at"`
}