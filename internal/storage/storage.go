package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	FilesDir  string
	ImagesDir string
}

func New(basePath string) *Storage {
	filesDir := filepath.Join(basePath, "files")
	imagesDir := filepath.Join(basePath, "images")
	
	os.MkdirAll(filesDir, 0755)
	os.MkdirAll(imagesDir, 0755)
	
	return &Storage{
		FilesDir:  filesDir,
		ImagesDir: imagesDir,
	}
}

func (s *Storage) SaveFile(filename string, content io.Reader) error {
	filePath := filepath.Join(s.FilesDir, filename)
	
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()
	
	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}
	
	return nil
}

func (s *Storage) GetFile(filename string) (*os.File, error) {
	filePath := filepath.Join(s.FilesDir, filename)
	return os.Open(filePath)
}

func (s *Storage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.FilesDir, filename)
	return os.Remove(filePath)
}

func (s *Storage) ListFiles() ([]string, error) {
	entries, err := os.ReadDir(s.FilesDir)
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	
	return files, nil
}

func (s *Storage) SaveImage(imageName string, content io.Reader) error {
	imagePath := filepath.Join(s.ImagesDir, imageName+".tar.gz")
	
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("create image file: %w", err)
	}
	defer file.Close()
	
	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("save image: %w", err)
	}
	
	return nil
}

func (s *Storage) GetImage(imageName string) (*os.File, error) {
	imagePath := filepath.Join(s.ImagesDir, imageName+".tar.gz")
	return os.Open(imagePath)
}

func (s *Storage) DeleteImage(imageName string) error {
	imagePath := filepath.Join(s.ImagesDir, imageName+".tar.gz")
	return os.Remove(imagePath)
}

func (s *Storage) ListImages() ([]string, error) {
	entries, err := os.ReadDir(s.ImagesDir)
	if err != nil {
		return nil, err
	}
	
	var images []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
			name := entry.Name()
			if len(name) > 7 && name[len(name)-7:] == ".tar.gz" {
				images = append(images, name[:len(name)-7])
			}
		}
	}
	
	return images, nil
}