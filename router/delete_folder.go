package main

import (
	"log"
	"os"
)

func deleteFolder(folderPath string) {
	if err := os.RemoveAll(folderPath); err != nil {
		log.Printf("Failed to delete %s: %v.\n", folderPath, err)
	}
}
