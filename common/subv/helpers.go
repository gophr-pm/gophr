package subv

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	// TODO figure out how to get ssh working
	//ret, cred := git.NewCredSshKey("git", "/Users/shikkic/.ssh/id_rsa.pub", "/Users/shikkic/.ssh/id_rsa", "")
	ret, cred := git.NewCredUserpassPlaintext("gophrpm", "PASSWORD_HERE")
	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

func deleteAchriveFile(archivePath string) error {
	log.Println(archivePath)
	if err := os.Remove(archivePath); err != nil {
		log.Println(err)
		return fmt.Errorf("Error, could not delete ref archive file. %v", err)
	}
	return nil
}

func deleteFolder(folderPath string) error {
	log.Println(folderPath)
	err := os.RemoveAll(folderPath)
	log.Println(err)
	return err
}

func unzip(archive, target string) error {
	log.Println("archive = ", archive)
	log.Println("target = ", target)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
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
