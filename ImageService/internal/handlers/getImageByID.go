package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"log"
	"main/internal/db"
	"net/http"
)

func GetImageByID(w http.ResponseWriter, r *http.Request) {
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

	downloadStream, err := bucket.OpenDownloadStream(objectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось открыть поток чтения файла из GridFS: %v", err), http.StatusInternalServerError)
		return
	}
	defer downloadStream.Close()

	file := downloadStream.GetFile()
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось получить информацию о файле из GridFS: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Length))

	if _, err := io.Copy(w, downloadStream); err != nil {
		http.Error(w, fmt.Sprintf("Не удалось отправить содержимое файла в ответ: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Файл успешно отправлен клиенту: %s", id)
}
