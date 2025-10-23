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

	registerHandler := &handlers.UserHandler{DB: conn}

	router.HandleFunc("/api/register", registerHandler.Register).Methods("POST")

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}