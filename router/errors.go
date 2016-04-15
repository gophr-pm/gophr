package main

import (
	"fmt"
	"net/http"
)

// TODO(skeswa): centralize all the errors here, give them codes, add logging

func respondWithInvalidURL(resp http.ResponseWriter, url string) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write([]byte(fmt.Sprintf(
		"Failed to process URL \"%s\". Please refer to the gophr docs for information on how to use it.",
		url,
	)))
}

func respondWithError(resp http.ResponseWriter, err error) {
	// TODO(skeswa): customize with custom formatting logic for different errors
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(err.Error()))
}

/************************** INVALID PACKAGE REQUEST ***************************/

// InvalidPackageRequestError is an error that occurs when an incoming package
// request could not be processed successfully.
type InvalidPackageRequestError struct {
	RequestURL string
	CausedBy   []error
}

// NewInvalidPackageRequestError creates a new InvalidPackageRequestError.
func NewInvalidPackageRequestError(
	requestURL string,
	causes ...error,
) InvalidPackageRequestError {
	return InvalidPackageRequestError{CausedBy: causes, RequestURL: requestURL}
}

func (err InvalidPackageRequestError) Error() string {
	return fmt.Sprintf(
		`Failed to parse and process request with URL "%s": %v`,
		err.RequestURL,
		err.Causes,
	)
}

// Causes returns the error(s) that caused this error.
func (err InvalidPackageRequestError) Causes() []error {
	return err.CausedBy
}

// PublicError returns an outside-friendly error message, and a
// corresponding status code.
func (err InvalidPackageRequestError) PublicError() (int, string) {
	return http.StatusBadRequest, fmt.Sprintf(
		`Failed to process the package request "%s". Please refer to the docs to find out how to use gophr.`,
		err.RequestURL,
	)
}

/********************** INVALID PACKAGE REF REQUEST ***********************/

// InvalidPackageRefRequestURLError is an error that occurs when an incoming
// request URL is an invalid package version request.
type InvalidPackageRefRequestURLError struct {
	RequestURL string
	CausedBy   []error
}

// NewInvalidPackageRefRequestURLError creates a new
// InvalidPackageRefRequestURLError.
func NewInvalidPackageRefRequestURLError(
	requestURL string,
	causes ...error,
) InvalidPackageRefRequestURLError {
	return InvalidPackageRefRequestURLError{
		RequestURL: requestURL,
		CausedBy:   causes,
	}
}

func (err InvalidPackageRefRequestURLError) Error() string {
	return fmt.Sprintf(
		`"%s" is not a valid package ref request URL.`,
		err.RequestURL,
	)
}

// Causes returns the error(s) that caused this error.
func (err InvalidPackageRefRequestURLError) Causes() []error {
	return err.CausedBy
}

/************************ INVALID BARE PACKAGE REQUEST ************************/

// InvalidBarePackageRequestURLError is an error that occurs when an incoming
// request URL is an invalid bare package request.
type InvalidBarePackageRequestURLError struct {
	RequestURL string
}

// NewInvalidBarePackageRequestURLError creates a new
// InvalidBarePackageRequestURLError.
func NewInvalidBarePackageRequestURLError(
	requestURL string,
) InvalidBarePackageRequestURLError {
	return InvalidBarePackageRequestURLError{RequestURL: requestURL}
}

func (err InvalidBarePackageRequestURLError) Error() string {
	return fmt.Sprintf(
		`"%s" is not a valid package version request URL.`,
		err.RequestURL,
	)
}

/********************** INVALID PACKAGE VERSION REQUEST ***********************/

// InvalidPackageVersionRequestURLError is an error that occurs when an incoming
// request URL is an invalid package version request.
type InvalidPackageVersionRequestURLError struct {
	RequestURL string
	CausedBy   []error
}

// NewInvalidPackageVersionRequestURLError creates a new
// InvalidPackageVersionRequestURLError.
func NewInvalidPackageVersionRequestURLError(
	requestURL string,
	causes ...error,
) InvalidPackageVersionRequestURLError {
	return InvalidPackageVersionRequestURLError{
		RequestURL: requestURL,
		CausedBy:   causes,
	}
}

func (err InvalidPackageVersionRequestURLError) Error() string {
	return fmt.Sprintf(
		`"%s" is not a valid package version request URL.`,
		err.RequestURL,
	)
}

// Causes returns the error(s) that caused this error.
func (err InvalidPackageVersionRequestURLError) Causes() []error {
	return err.CausedBy
}

/********************** NO SUCH PACKAGE VERSION REQUEST ***********************/

// NoSuchPackageVersionError is an error that occurs when a equested package
// version doesn't exist.
type NoSuchPackageVersionError struct {
	PackageAuthor string
	PackageRepo   string
	Selector      string
	CausedBy      []error
}

// NewNoSuchPackageVersionError creates a new NoSuchPackageVersionError.
func NewNoSuchPackageVersionError(
	packageAuthor string,
	packageRepo string,
	selector string,
	causes ...error,
) NoSuchPackageVersionError {
	return NoSuchPackageVersionError{
		PackageAuthor: packageAuthor,
		PackageRepo:   packageRepo,
		Selector:      selector,
		CausedBy:      causes,
	}
}

func (err NoSuchPackageVersionError) Error() string {
	return fmt.Sprintf(
		`Could not find a version of "%s/%s" that matches "%s".`,
		err.PackageAuthor,
		err.PackageRepo,
		err.Selector,
	)
}

// Causes returns the error(s) that caused this error.
func (err NoSuchPackageVersionError) Causes() []error {
	return err.CausedBy
}

// PublicError returns an outside-friendly error message, and a
// corresponding status code.
func (err NoSuchPackageVersionError) PublicError() (int, string) {
	return http.StatusNotFound, err.Error()
}
