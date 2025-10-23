package models

import "time"

type Book struct {
    ID int `db:"id" json:"id"`
    Title string `db:"title" json:"title"`
    Author string `db:"author" json:"author"`
    CategoryID *int `db:"category_id" json:"category_id"`
    Stock int `db:"stock" json:"stock"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}