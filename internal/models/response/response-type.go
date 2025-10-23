package response

import "time"

type BookResponse struct {
	ID         int       `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Author     string    `db:"author" json:"author"`
	CategoryID *int      `db:"category_id" json:"category_id"`
	Category   *string   `db:"category" json:"category"`
	Stock      int       `db:"stock" json:"stock"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
type UserBorrowingResponse struct {
	ID         int        `db:"id" json:"id"`
	BookID     int        `db:"book_id" json:"book_id"`
	BookTitle  string     `db:"book_title" json:"book_title"`
	Author     string     `db:"author" json:"author"`
	BorrowedAt time.Time  `db:"borrowed_at" json:"borrowed_at"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at"` // pakai pointer karena bisa NULL
	Status     string     `db:"status" json:"status"`
}
