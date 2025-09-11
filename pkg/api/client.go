package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dollarkillerx/unregistry/pkg/progress"
)

type Client struct {
	BaseURL string
	Token   string
	client  *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		client:  &http.Client{},
	}
}

func (c *Client) doRequest(method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.BaseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.client.Do(req)
}

// File operations
func (c *Client) UploadFile(filePath string) error {
	return c.UploadFileWithProgress(filePath, false)
}

func (c *Client) UploadFileWithProgress(filePath string, showProgress bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get file info: %w", err)
	}

	// Create a pipe for streaming
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	contentType := writer.FormDataContentType()

	// Start a goroutine to write to the pipe
	go func() {
		defer pipeWriter.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}
	}()

	var body io.Reader = pipeReader
	if showProgress {
		// Calculate the approximate size of the multipart request
		// This is an estimation: file size + multipart overhead (roughly 200-300 bytes)
		totalSize := fileInfo.Size() + 300
		progressReader := progress.NewReader(pipeReader, totalSize, "Uploading "+filepath.Base(filePath))
		defer progressReader.Close()
		body = progressReader
	}

	resp, err := c.doRequest("POST", "/api/file/upload", body, contentType)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(respBody))
	}

	return nil
}

func (c *Client) DownloadFile(filename, destPath string) error {
	return c.DownloadFileWithProgress(filename, destPath, false)
}

func (c *Client) DownloadFileWithProgress(filename, destPath string, showProgress bool) error {
	resp, err := c.doRequest("GET", "/api/file/download/"+filename, nil, "")
	if err != nil {
		return fmt.Errorf("download request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed: %s", string(body))
	}

	if destPath == "" {
		destPath = filename
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer file.Close()

	if showProgress {
		contentLength := resp.Header.Get("Content-Length")
		if contentLength != "" {
			size, _ := strconv.ParseInt(contentLength, 10, 64)
			progressReader := progress.NewReader(resp.Body, size, "Downloading "+filename)
			defer progressReader.Close()
			_, err = io.Copy(file, progressReader)
		} else {
			_, err = io.Copy(file, resp.Body)
		}
	} else {
		_, err = io.Copy(file, resp.Body)
	}
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (c *Client) ListFiles() ([]string, error) {
	resp, err := c.doRequest("GET", "/api/file/list", nil, "")
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed: %s", string(body))
	}

	var result struct {
		Files []string `json:"files"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Files, nil
}

func (c *Client) DeleteFile(filename string) error {
	resp, err := c.doRequest("DELETE", "/api/file/"+filename, nil, "")
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(body))
	}

	return nil
}

// Image operations
func (c *Client) UploadImage(imagePath string) error {
	return c.UploadImageWithProgress(imagePath, false)
}

func (c *Client) UploadImageWithProgress(imagePath string, showProgress bool) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get image info: %w", err)
	}

	// Create a pipe for streaming
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	contentType := writer.FormDataContentType()

	// Start a goroutine to write to the pipe
	go func() {
		defer pipeWriter.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}
	}()

	var body io.Reader = pipeReader
	if showProgress {
		// Calculate the approximate size of the multipart request
		// This is an estimation: file size + multipart overhead (roughly 200-300 bytes)
		totalSize := fileInfo.Size() + 300
		progressReader := progress.NewReader(pipeReader, totalSize, "Uploading "+filepath.Base(imagePath))
		defer progressReader.Close()
		body = progressReader
	}

	resp, err := c.doRequest("POST", "/api/img/upload", body, contentType)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(respBody))
	}

	return nil
}

func (c *Client) DownloadImage(imageName, destPath string) error {
	return c.DownloadImageWithProgress(imageName, destPath, false)
}

func (c *Client) DownloadImageWithProgress(imageName, destPath string, showProgress bool) error {
	resp, err := c.doRequest("GET", "/api/img/download/"+imageName, nil, "")
	if err != nil {
		return fmt.Errorf("download request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed: %s", string(body))
	}

	if destPath == "" {
		destPath = imageName + ".tar.gz"
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer file.Close()

	if showProgress {
		contentLength := resp.Header.Get("Content-Length")
		if contentLength != "" {
			size, _ := strconv.ParseInt(contentLength, 10, 64)
			progressReader := progress.NewReader(resp.Body, size, "Downloading "+imageName)
			defer progressReader.Close()
			_, err = io.Copy(file, progressReader)
		} else {
			_, err = io.Copy(file, resp.Body)
		}
	} else {
		_, err = io.Copy(file, resp.Body)
	}
	if err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	return nil
}

func (c *Client) ListImages() ([]string, error) {
	resp, err := c.doRequest("GET", "/api/img/list", nil, "")
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed: %s", string(body))
	}

	var result struct {
		Images []string `json:"images"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Images, nil
}

func (c *Client) DeleteImage(imageName string) error {
	resp, err := c.doRequest("DELETE", "/api/img/"+imageName, nil, "")
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(body))
	}

	return nil
}