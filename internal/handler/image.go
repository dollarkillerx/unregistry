package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dollarkillerx/unregistry/internal/storage"
	"github.com/gorilla/mux"
)

type ImageHandler struct {
	storage *storage.Storage
}

func NewImageHandler(storage *storage.Storage) *ImageHandler {
	return &ImageHandler{storage: storage}
}

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extract image name from filename (remove .tar.gz extension if present)
	imageName := header.Filename
	if len(imageName) > 7 && imageName[len(imageName)-7:] == ".tar.gz" {
		imageName = imageName[:len(imageName)-7]
	}

	err = h.storage.SaveImage(imageName, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Image uploaded successfully",
		"image_name": imageName,
	})
}

func (h *ImageHandler) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["name"]

	file, err := h.storage.GetImage(imageName)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar.gz", imageName))
	
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to download image", http.StatusInternalServerError)
		return
	}
}

func (h *ImageHandler) List(w http.ResponseWriter, r *http.Request) {
	images, err := h.storage.ListImages()
	if err != nil {
		http.Error(w, "Failed to list images", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images": images,
	})
}

func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["name"]

	err := h.storage.DeleteImage(imageName)
	if err != nil {
		http.Error(w, "Failed to delete image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image deleted successfully",
	})
}