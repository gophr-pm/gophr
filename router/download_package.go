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

func downloadPackage(args packageDownloaderArgs) (packageDownloadPaths, error) {
	downloadPaths := packageDownloadPaths{}

	workDirPath := filepath.Join(args.constructionZonePath, generateWorkDirName())
	if err := os.Mkdir(workDirPath, 0644); err != nil {
		return downloadPaths, fmt.Errorf("Could not create workDir %s: %v.", workDirPath, err)
	}

	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", args.author, args.repo, args.sha)
	zipResp, err := http.Get(zipURL)
	defer zipResp.Body.Close()
	if err != nil || zipResp.StatusCode == 404 {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not find args.sha archive for %s: %v.", zipURL, err)
	}

	zipFilePath := filepath.Join(workDirPath, "archive.zip")
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

	if err = unzip(zipFilePath, workDirPath); err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could not unzip to file system: %v.", err)
	}

	// Finding the name of the archive folder.
	files, err := ioutil.ReadDir(workDirPath)
	if err != nil {
		defer deleteFolder(workDirPath)
		return downloadPaths, fmt.Errorf("Could look up files in workDir: %v.", err)
	}
	for _, f := range files {
		fileName := f.Name()
		if fileName != "archive.zip" {
			downloadPaths.workDirPath = workDirPath
			downloadPaths.archiveDirPath = filepath.Join(workDirPath, fileName)
			return downloadPaths, nil
		}
	}

	return downloadPaths, errors.New("Could not find archiveDirPath.")
}

func unzip(archive, target string) error {
	log.Println("archive = ", archive)
	log.Println("target = ", target)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0644); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
