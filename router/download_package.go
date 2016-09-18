package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	packageZipFileName   = "archive.zip"
	githubZipURLTemplate = "https://github.com/%s/%s/archive/%s.zip"
)

// downloadPackage downloads a go package repository from Github into the
// construction zone and returns the created directories.
func downloadPackage(args packageDownloaderArgs) (packageDownloadPaths, error) {
	downloadPaths := packageDownloadPaths{}

	// Create the working directory. Everything else in this functions happens in
	// here.
	workDirPath := filepath.Join(args.constructionZonePath, generateWorkDirName())
	if err := os.Mkdir(workDirPath, 0644); err != nil {
		return downloadPaths, fmt.Errorf("Could not create workDir %s: %v.", workDirPath, err)
	}

	// Use a zip strategy to download from Github in order to save of data
	// transfer and on-disk storage needs.
	zipURL := fmt.Sprintf(githubZipURLTemplate, args.author, args.repo, args.sha)
	zipResp, err := http.Get(zipURL)
	defer zipResp.Body.Close()
	if err != nil || zipResp.StatusCode == 404 {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not find args.sha archive for %s: %v.", zipURL, err)
	}

	// Create the zip file and unzip it.
	zipFilePath := filepath.Join(workDirPath, packageZipFileName)
	zipFile, err := os.Create(zipFilePath)
	defer zipFile.Close()
	if err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not write archive to file system: %v.", err)
	}
	if _, err = io.Copy(zipFile, zipResp.Body); err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not copy archive to file system: %v.", err)
	}
	if err = unzipPackageZip(zipFilePath, workDirPath); err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not unzip to file system: %v.", err)
	}

	// Find the name of the archive folder so that we can return it.
	files, err := ioutil.ReadDir(workDirPath)
	if err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could look up files in workDir: %v.", err)
	}
	for _, f := range files {
		fileName := f.Name()
		if fileName != packageZipFileName {
			downloadPaths.workDirPath = workDirPath
			downloadPaths.archiveDirPath = filepath.Join(workDirPath, fileName)
			return downloadPaths, nil
		}
	}

	return downloadPaths, errors.New("Could not find archiveDirPath.")
}

// unzipPackageZip unzips a package zip (called archive) into the target
// directory.
func unzipPackageZip(archive, target string) error {
	// TODO(skeswa): get rid of noisy logs.
	log.Println("archive = ", archive)
	log.Println("target = ", target)
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
