package main

import (
	"log"
	"net/http"
	"os"

	"github.com/faqq11/lib-management/internal/db"
	"github.com/faqq11/lib-management/internal/handlers"
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
	bookHandler := &handlers.BookHandler{DB:conn}
	categoryHandler := &handlers.CategoryHandler{DB:conn}

	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	router.HandleFunc("/api/create-book", bookHandler.InsertBook).Methods("POST")
	router.HandleFunc("/api/books", bookHandler.GetAllBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", bookHandler.GetBookById).Methods("GET")
	router.HandleFunc("/api/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}/increase-stock", bookHandler.IncreaseStock).Methods("PUT")
	router.HandleFunc("/api/books/{id}/decrease-stock", bookHandler.DecreaseStock).Methods("PUT")
	router.HandleFunc("/api/books/{id}/delete", bookHandler.DeleteBook).Methods("DELETE")

	router.HandleFunc("/api/create-category", categoryHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/api/delete-category/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}