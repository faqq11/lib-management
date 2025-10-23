package main

import (
	"log"
	"net/http"
	"os"

	"github.com/faqq11/lib-management/internal/db"
	"github.com/faqq11/lib-management/internal/handlers"
	"github.com/faqq11/lib-management/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	conn, err := db.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	router := mux.NewRouter()

	userHandler := &handlers.UserHandler{DB: conn}
	bookHandler := &handlers.BookHandler{DB: conn}
	categoryHandler := &handlers.CategoryHandler{DB: conn}
	borrowHandler := &handlers.BorrowHandler{DB: conn}

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	adminOnly := protected.PathPrefix("").Subrouter()
	adminOnly.Use(middleware.AdminOnlyMiddleware)

	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	adminOnly.HandleFunc("/create-book", bookHandler.InsertBook).Methods("POST")
	protected.HandleFunc("/books", bookHandler.GetAllBooks).Methods("GET")
	protected.HandleFunc("/books/search", bookHandler.SearchBooks).Methods("GET")
	protected.HandleFunc("/books/{id}", bookHandler.GetBookById).Methods("GET")
	adminOnly.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	adminOnly.HandleFunc("/books/{id}/increase-stock", bookHandler.IncreaseStock).Methods("PUT")
	adminOnly.HandleFunc("/books/{id}/decrease-stock", bookHandler.DecreaseStock).Methods("PUT")
	adminOnly.HandleFunc("/books/{id}/delete", bookHandler.DeleteBook).Methods("DELETE")

	adminOnly.HandleFunc("/create-category", categoryHandler.CreateCategory).Methods("POST")
	adminOnly.HandleFunc("/delete-category/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	protected.HandleFunc("/my-borrowings", borrowHandler.GetUserBorrowings).Methods("GET")
	protected.HandleFunc("/books/{id}/borrow", borrowHandler.BorrowBook).Methods("POST")
	protected.HandleFunc("/borrowings/{id}/return", borrowHandler.ReturnBook).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
