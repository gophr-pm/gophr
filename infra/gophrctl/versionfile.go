package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type imageVersionPart int

const (
	imageVersionMajorPart = imageVersionPart(1)
	imageVersionMinorPart = imageVersionPart(2)
	imageVersionPatchPart = imageVersionPart(3)
)

type imageVersion struct {
	major int
	minor int
	patch int
}

func (iv imageVersion) buffer() *bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(iv.major))
	buffer.WriteByte('.')
	buffer.WriteString(strconv.Itoa(iv.minor))
	buffer.WriteByte('.')
	buffer.WriteString(strconv.Itoa(iv.patch))
	return &buffer
}

func (iv imageVersion) Bytes() []byte {
	return iv.buffer().Bytes()
}

func (iv imageVersion) String() string {
	return iv.buffer().String()
}

func newInvalidVersionFileError(versionFilePath, versionFileStr string) error {
	return fmt.Errorf("Invalid version in file \"%s\": %s", versionFilePath, versionFileStr)
}

func readImageVersion(versionFilePath string) (imageVersion, error) {
	version := imageVersion{}
	versionFileData, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		return version, err
	}

	versionFileStr := string(versionFileData[:])
	versionFileStr = strings.TrimSpace(versionFileStr)
	versionPartStrs := strings.Split(versionFileStr, ".")

	if len(versionPartStrs) < 3 {
		return version, newInvalidVersionFileError(versionFilePath, versionFileStr)
	}

	if version.major, err = strconv.Atoi(versionPartStrs[0]); err != nil {
		return version, newInvalidVersionFileError(versionFilePath, versionFileStr)
	} else if version.minor, err = strconv.Atoi(versionPartStrs[1]); err != nil {
		return version, newInvalidVersionFileError(versionFilePath, versionFileStr)
	} else if version.patch, err = strconv.Atoi(versionPartStrs[2]); err != nil {
		return version, newInvalidVersionFileError(versionFilePath, versionFileStr)
	}

	return version, nil
}

func bumpImageVersion(version imageVersion, part imageVersionPart) imageVersion {
	switch part {
	case imageVersionMajorPart:
		return imageVersion{
			major: version.major + 1,
			minor: 0,
			patch: 0,
		}
	case imageVersionMinorPart:
		return imageVersion{
			major: version.major,
			minor: version.minor + 1,
			patch: 0,
		}
	case imageVersionPatchPart:
		return imageVersion{
			major: version.major,
			minor: version.minor,
			patch: version.patch + 1,
		}
	}

	return version
}

func writeToVersionFile(versionFilePath string, version imageVersion) error {
	return ioutil.WriteFile(versionFilePath, version.Bytes(), 0644)
}

func promptImageVersionBump(imageName, versionFilePath string) (imageVersion, error) {
	version, err := readImageVersion(versionFilePath)
	if err != nil {
		return version, err
	}

	fmt.Print("The current image version of ")
	blue.Print(imageName)
	fmt.Print(" is ")
	blue.Print(version.String())
	fmt.Println(".\nWhat level of this version should be bumped?")
	yellow.Print("(0) ")
	fmt.Println("None")
	yellow.Print("(1) ")
	fmt.Println("Major")
	yellow.Print("(2) ")
	fmt.Println("Minor")
	yellow.Print("(3) ")
	fmt.Println("Patch")
	fmt.Println()

	var (
		level         int
		scanner       = bufio.NewScanner(os.Stdin)
		bumpedVersion imageVersion
	)

	for {
		fmt.Print("Selected Level [3]: ")
		scanner.Scan()
		text := scanner.Text()

		if len(text) == 0 {
			bumpedVersion = bumpImageVersion(version, imageVersionPatchPart)
			break
		} else if level, err = strconv.Atoi(text); err != nil || level < 0 || level > 3 {
			printError("Invalid bump level chosen. Try again.")
		} else {
			if level == 0 {
				printSuccess("Version was left unbumped\n")
				return version, nil
			} else if level == 1 {
				bumpedVersion = bumpImageVersion(version, imageVersionMajorPart)
				break
			} else if level == 2 {
				bumpedVersion = bumpImageVersion(version, imageVersionMinorPart)
				break
			} else if level == 3 {
				bumpedVersion = bumpImageVersion(version, imageVersionPatchPart)
				break
			}
		}
	}

	// Write to the version file or bust.
	if err = writeToVersionFile(versionFilePath, bumpedVersion); err != nil {
		return bumpedVersion, err
	}

	// Print the success message in pieces due to the multiple colors present.
	green.Print("✓ Image version was bumped from ")
	blue.Print(version.String())
	green.Print(" to ")
	blue.Println(bumpedVersion.String())
	fmt.Println()

	return bumpedVersion, nil
}
