package main

import (
	"github.com/gorilla/mux"
	"main/internal/db"
	"main/internal/handlers"
	"net/http"
)

func main() {
	db.ConnectDB()
	router := mux.NewRouter()
	router.HandleFunc("/uploadImage", handlers.UploadImage).Methods("POST")
	router.HandleFunc("/getImageByID/{id}", handlers.GetImageByID).Methods("GET")
	router.HandleFunc("/deleteImageByID/{id}", handlers.DeleteImageByID).Methods("DELETE")
	http.ListenAndServe("localhost:8080", router)
}
