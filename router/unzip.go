package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// unzipArchive unzips a zip archive into the target directory.
func unzipArchive(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	// Use the zip reader to identify and create files in the filesystem from
	// the zip.
	for _, file := range reader.File {
		// If the file is a directory, make sure its full path exists.
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// Now we know that file is a File. Lets open it so as to copy it into the
		// filesystem.
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		// Get the file descriptor for file.
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// Use the file descriptor to perform a copy.
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
