package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dollarkillerx/unregistry/pkg/api"
	"github.com/dollarkillerx/unregistry/pkg/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "unrg",
	Short: "Unregistry client - A private file/image storage system",
	Long:  "Command line client for Unregistry private file and Docker image storage system.",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var setTokenCmd = &cobra.Command{
	Use:   "set-token <token>",
	Short: "Set authentication token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		
		cfg.Token = token
		
		err = cfg.Save()
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Token saved successfully")
	},
}

var setURLCmd = &cobra.Command{
	Use:   "set-url <url>",
	Short: "Set server base URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		
		cfg.BaseURL = url
		
		err = cfg.Save()
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Base URL set to: %s\n", url)
	},
}

// File commands
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File operations",
}

var filePushCmd = &cobra.Command{
	Use:   "push <filepath>",
	Short: "Upload a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		err = client.UploadFileWithProgress(filePath, true)
		if err != nil {
			fmt.Printf("Upload failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("\nFile %s uploaded successfully\n", filepath.Base(filePath))
	},
}

var filePullCmd = &cobra.Command{
	Use:   "pull <filename> [dest]",
	Short: "Download a file",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		var dest string
		if len(args) > 1 {
			dest = args[1]
		}
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		err = client.DownloadFileWithProgress(filename, dest, true)
		if err != nil {
			fmt.Printf("Download failed: %v\n", err)
			os.Exit(1)
		}
		
		if dest == "" {
			dest = filename
		}
		fmt.Printf("\nFile downloaded to: %s\n", dest)
	},
}

var fileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all files",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		files, err := client.ListFiles()
		if err != nil {
			fmt.Printf("List failed: %v\n", err)
			os.Exit(1)
		}
		
		if len(files) == 0 {
			fmt.Println("No files found")
			return
		}
		
		fmt.Println("Files:")
		for _, file := range files {
			fmt.Printf("  %s\n", file)
		}
	},
}

var fileDeleteCmd = &cobra.Command{
	Use:   "delete <filename>",
	Short: "Delete a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		err = client.DeleteFile(filename)
		if err != nil {
			fmt.Printf("Delete failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("File %s deleted successfully\n", filename)
	},
}

// Image commands
var imgCmd = &cobra.Command{
	Use:   "img",
	Short: "Docker image operations",
}

var imgPushCmd = &cobra.Command{
	Use:   "push <docker_image>",
	Short: "Push a Docker image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dockerImage := args[0]
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		// Create temporary file for docker save
		tempFile := fmt.Sprintf("/tmp/%s.tar.gz", strings.ReplaceAll(dockerImage, "/", "_"))
		defer os.Remove(tempFile)
		
		fmt.Println("Preparing Docker image...")
		saveCmd := exec.Command("docker", "save", dockerImage)
		gzipCmd := exec.Command("gzip")
		
		// Pipe docker save output to gzip
		gzipCmd.Stdin, _ = saveCmd.StdoutPipe()
		
		// Create output file
		outFile, err := os.Create(tempFile)
		if err != nil {
			fmt.Printf("Failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()
		
		gzipCmd.Stdout = outFile
		
		// Start gzip first, then docker save
		if err := gzipCmd.Start(); err != nil {
			fmt.Printf("Failed to start gzip: %v\n", err)
			os.Exit(1)
		}
		
		if err := saveCmd.Run(); err != nil {
			fmt.Printf("Docker save failed: %v\n", err)
			os.Exit(1)
		}
		
		if err := gzipCmd.Wait(); err != nil {
			fmt.Printf("Gzip failed: %v\n", err)
			os.Exit(1)
		}
		
		err = client.UploadImageWithProgress(tempFile, true)
		if err != nil {
			fmt.Printf("Upload failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("\nImage %s pushed successfully\n", dockerImage)
	},
}

var imgPullCmd = &cobra.Command{
	Use:   "pull <docker_image>",
	Short: "Pull a Docker image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dockerImage := args[0]
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		// Create temporary file
		tempFile := fmt.Sprintf("/tmp/%s.tar.gz", strings.ReplaceAll(dockerImage, "/", "_"))
		defer os.Remove(tempFile)
		
		err = client.DownloadImageWithProgress(dockerImage, tempFile, true)
		if err != nil {
			fmt.Printf("Download failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Loading Docker image...")
		loadCmd := exec.Command("docker", "load")
		
		// Open compressed file and decompress
		gzipCmd := exec.Command("gunzip", "-c", tempFile)
		loadCmd.Stdin, _ = gzipCmd.StdoutPipe()
		
		if err := loadCmd.Start(); err != nil {
			fmt.Printf("Failed to start docker load: %v\n", err)
			os.Exit(1)
		}
		
		if err := gzipCmd.Run(); err != nil {
			fmt.Printf("Gunzip failed: %v\n", err)
			os.Exit(1)
		}
		
		if err := loadCmd.Wait(); err != nil {
			fmt.Printf("Docker load failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("\nImage %s pulled successfully\n", dockerImage)
	},
}

var imgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all images",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		images, err := client.ListImages()
		if err != nil {
			fmt.Printf("List failed: %v\n", err)
			os.Exit(1)
		}
		
		if len(images) == 0 {
			fmt.Println("No images found")
			return
		}
		
		fmt.Println("Images:")
		for _, image := range images {
			fmt.Printf("  %s\n", image)
		}
	},
}

var imgDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		
		client, err := getClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		err = client.DeleteImage(imageName)
		if err != nil {
			fmt.Printf("Delete failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Image %s deleted successfully\n", imageName)
	},
}

func getClient() (*api.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	
	if cfg.Token == "" {
		return nil, fmt.Errorf("no token configured. Use 'unrg config set-token <token>' first")
	}
	
	return api.NewClient(cfg.BaseURL, cfg.Token), nil
}

func init() {
	// Config commands
	configCmd.AddCommand(setTokenCmd)
	configCmd.AddCommand(setURLCmd)
	rootCmd.AddCommand(configCmd)
	
	// File commands
	fileCmd.AddCommand(filePushCmd)
	fileCmd.AddCommand(filePullCmd)
	fileCmd.AddCommand(fileListCmd)
	fileCmd.AddCommand(fileDeleteCmd)
	rootCmd.AddCommand(fileCmd)
	
	// Image commands
	imgCmd.AddCommand(imgPushCmd)
	imgCmd.AddCommand(imgPullCmd)
	imgCmd.AddCommand(imgListCmd)
	imgCmd.AddCommand(imgDeleteCmd)
	rootCmd.AddCommand(imgCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}