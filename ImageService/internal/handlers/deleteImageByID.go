package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"log"
	"main/internal/db"
	"net/http"
)

func DeleteImageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неверный формат ID файла: %v", err), http.StatusBadRequest)
		return
	}

	bucket, err := gridfs.NewBucket(
		db.MongoClient.Database("images"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось подключиться к GridFS: %v", err), http.StatusInternalServerError)
		return
	}

	err = bucket.Delete(objectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось удалить файл из GridFS: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Файл с ID %s успешно удален из GridFS", id)

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Файл с ID %s успешно удален", id)
}
