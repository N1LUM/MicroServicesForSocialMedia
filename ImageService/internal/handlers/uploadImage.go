package handlers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"log"
	"main/internal/db"
	"net/http"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Не удалось преоразовать картинку для чтения: %v", err), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось получить файл из запроса: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Загружен файл: %s", handler.Filename)
	log.Printf("Размер файла: %d", handler.Size)
	log.Printf("MIME тип файла: %s", handler.Header.Get("Content-Type"))

	// Подключение к GridFS
	bucket, err := gridfs.NewBucket(
		db.MongoClient.Database("images"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось подключиться к GridFS: %v", err), http.StatusInternalServerError)
		return
	}

	// Создание загрузочного стрима в GridFS
	uploadStream, err := bucket.OpenUploadStream(handler.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось создать загрузочный поток в GridFS: %v", err), http.StatusInternalServerError)
		return
	}
	defer uploadStream.Close()

	// Сохранение файла в GridFS
	fileSize, err := io.Copy(uploadStream, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось сохранить файл в GridFS: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Размер загруженного файла в GridFS: %d", fileSize)

	// Получение ID загруженного файла
	objectID := uploadStream.FileID
	log.Printf("ID загруженного файла: %s", objectID)

	// Возвращение успешного ответа клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл успешно загружен, ID: %s", objectID)
}
