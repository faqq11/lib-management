package models

import "time"

type User struct {
    ID int `db:"id" json:"id"`
    Username string `db:"username" json:"username"`
    Password string `db:"password,omitempty" json:"-"`
    Role string `db:"role" json:"role"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}