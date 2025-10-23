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