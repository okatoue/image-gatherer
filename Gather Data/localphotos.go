package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// isImageOrVideo checks if a file is an image or video based on its extension
func isImageOrVideo(fileName string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	videoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm"}

	extension := strings.ToLower(filepath.Ext(fileName))
	for _, ext := range imageExtensions {
		if extension == ext {
			return true
		}
	}
	for _, ext := range videoExtensions {
		if extension == ext {
			return true
		}
	}
	return false
}

// copyFile copies a file from the source path to the destination path
func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// checkFolderForMedia goes through a folder and copies any images or videos found to another folder
func checkFolderForMedia(folderPath, destFolderPath string) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageOrVideo(info.Name()) {
			destPath := filepath.Join(destFolderPath, info.Name())
			err := copyFile(path, destPath)
			if err != nil {
				fmt.Println("Error copying file:", err)
			} else {
				fmt.Println("Copied media file to:", destPath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

func main() {
	folderPath := "path to copy from goes here"
	destFolderPath := "Destination path goes here"

	// Create the destination folder if it doesn't exist
	if _, err := os.Stat(destFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(destFolderPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating destination folder:", err)
			return
		}
	}

	checkFolderForMedia(folderPath, destFolderPath)
}
