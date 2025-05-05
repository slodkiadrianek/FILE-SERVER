package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var uploadDir string = os.Getenv("UPLOAD_DIR")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	r := mux.NewRouter()
<<<<<<< HEAD
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(corsMiddleware)
=======
>>>>>>> 66bfa466bb1142f18d19f81b8a01d7ad29fe7dc9
	r.HandleFunc("/upload", uploadHandler).Methods("POST")

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://192.168.0.113:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(524288000)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var photoUrls []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}

		timestamp := time.Now().UnixNano()
		safeName := strings.ReplaceAll(fileHeader.Filename, " ", "-")
		newFilename := fmt.Sprintf("%d-%s", timestamp, safeName)
		dstPath := filepath.Join(uploadDir, newFilename)

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("%s://%s/uploads/%s", getScheme(r), r.Host, newFilename)
		photoUrls = append(photoUrls, url)
		file.Close()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message":true,"count":%d,"photoUrls":%q}`, len(photoUrls), photoUrls)))
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
