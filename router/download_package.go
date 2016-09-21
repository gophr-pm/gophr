package main

import (
	"errors"
	"fmt"
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
	if err := args.io.Mkdir(workDirPath, 0644); err != nil {
		return downloadPaths, fmt.Errorf("Could not create workDir %s: %v.", workDirPath, err)
	}

	// Use a zip strategy to download from Github in order to save of data
	// transfer and on-disk storage needs.
	zipURL := fmt.Sprintf(githubZipURLTemplate, args.author, args.repo, args.sha)
	zipResp, err := args.doHTTPGet(zipURL)
	defer zipResp.Body.Close()
	if err != nil || zipResp.StatusCode == 404 {
		defer args.deleteWorkDir(workDirPath)
		return downloadPaths, fmt.Errorf("Could not find args.sha archive for %s: %v.", zipURL, err)
	}

	// Create the zip file and unzip it.
	zipFilePath := filepath.Join(workDirPath, packageZipFileName)
	zipFile, err := args.io.Create(zipFilePath)
	defer zipFile.Close()
	if err != nil {
		args.deleteWorkDir(workDirPath)
		return downloadPaths, fmt.Errorf("Could not write archive to file system: %v.", err)
	}
	if _, err = args.io.Copy(zipFile, zipResp.Body); err != nil {
		args.deleteWorkDir(workDirPath)
		return downloadPaths, fmt.Errorf("Could not copy archive to file system: %v.", err)
	}
	if err = args.unzipArchive(zipFilePath, workDirPath); err != nil {
		args.deleteWorkDir(workDirPath)
		return downloadPaths, fmt.Errorf("Could not unzip to file system: %v.", err)
	}

	// Find the name of the archive folder so that we can return it.
	files, err := args.io.ReadDir(workDirPath)
	if err != nil {
		args.deleteWorkDir(workDirPath)
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
