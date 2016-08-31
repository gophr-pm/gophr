package subv

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// Exit Hook to clean up files
func cleanUpExitHook(folderName string) {
	log.Printf("Exit hook initiated deleting folder %s \n", folderName)
	if strings.HasPrefix(folderName, "/") == true {
		cmd := fmt.Sprintf("cd /tmp && rm -rf %s", folderName)
		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Println("Could not properly engage exit hook")
			log.Fatalln(err)
		}
		log.Printf("Output: %s", out)
	} else {
		log.Println("Cowardly refusing to initiate exit hook. Will not rm -rf folder names that contains any leading '/'")
	}
}

// Check Error and Engage Exit Hook if fatal error occured
func checkError(err error, folderName string) {
	if err != nil {
		cleanUpExitHook(folderName)
	}
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
