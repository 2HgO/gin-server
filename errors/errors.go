package errors

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
)

type ErrorType string

type ConstError string

func (c ConstError) Error() string {
	return string(c)
}

const (
	ErrNotFound       ErrorType = "ENTRY_NOT_FOUND"
	ErrValidation     ErrorType = "VALIDATION_ERROR"
	ErrEntryExists    ErrorType = "ENTRY_EXISTS"
	ErrEntryDeleted   ErrorType = "ENTRY_DELETED"
	ErrAuthorization  ErrorType = "AUTHORIZATION_ERROR"
	ErrExpiredToken   ErrorType = "TOKEN_EXPIRED"
	ErrAuthentication ErrorType = "AUTHENTICATION_ERROR"
	ErrInvalidToken   ErrorType = "TOKEN_INVALID"
	ErrPermission     ErrorType = "PERMISSION_ERROR"
	ErrFatal          ErrorType = "FATAL_ERROR"
	ErrNotImplemented ErrorType = "NOT_IMPLEMENTED_ERROR"
)

const (
	NotFoundError         ConstError = "Entry not found"
	ValidationError       ConstError = "Invalid request parameter"
	EntryExistsError      ConstError = "Entry already exists"
	EntryDeletedError     ConstError = "Entry has been deleted or deactivated"
	AuthorizationError    ConstError = "You are not authorized to view this resource"
	UnsupportedMediaError ConstError = "Invalid media/file type"
	ExpiredTokenError     ConstError = "Expired token"
	InvalidTokenError     ConstError = "Invalid token"
	AuthenticationError   ConstError = "User could not be authenticated"
	PermissionError       ConstError = "You do not have permission to perform this action"
	FatalError            ConstError = "An error has occured on our end"
	NotImplementedError   ConstError = "Handler or method has not been implemented"

	// internal error for db package
	ErrInvalidDoc ConstError = "Invalid Document"
)

func ErrorHandler(err error) (code int, message string, errType ErrorType) {
	if v, ok := err.(validator.ValidationErrors); ok {
		message = fmt.Sprintf("Validation failed on field { %s }, Condition: %s", v[0].Field(), v[0].ActualTag())
		if v[0].Param() != "" {
			message += fmt.Sprintf("{ %s }", v[0].Param())
		}
		if v[0].Value() != "" {
			message += fmt.Sprintf(", Value Recieved: %s", v[0].Value())
		}
		code, errType = http.StatusBadRequest, ErrValidation

		return code, message, errType
	}
	if _, ok := err.(mongo.WriteException); ok {
		return http.StatusConflict, EntryExistsError.Error(), ErrEntryExists
	}

	switch {
	case errors.Is(err, io.EOF):
		code, message, errType = http.StatusBadRequest, "No request body", ErrValidation
	case errors.Is(err, mongo.ErrNoDocuments), errors.Is(err, NotFoundError):
		code, message, errType = http.StatusNotFound, err.Error(), ErrNotFound
	case errors.Is(err, ValidationError):
		code, message, errType = http.StatusBadRequest, err.Error(), ErrValidation
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		code, message, errType = http.StatusUnauthorized, fmt.Sprintf("%s: Invalid login credentials", AuthenticationError.Error()), ErrAuthentication
	case errors.Is(err, InvalidTokenError):
		code, message, errType = http.StatusUnauthorized, err.Error(), ErrInvalidToken
	case errors.Is(err, ExpiredTokenError):
		code, message, errType = http.StatusUnauthorized, err.Error(), ErrExpiredToken
	case errors.Is(err, EntryDeletedError):
		code, message, errType = http.StatusNotFound, err.Error(), ErrEntryDeleted
	case errors.Is(err, EntryExistsError):
		code, message, errType = http.StatusConflict, err.Error(), ErrEntryExists
	case errors.Is(err, PermissionError):
		code, message, errType = http.StatusForbidden, err.Error(), ErrPermission
	case errors.Is(err, UnsupportedMediaError):
		code, message, errType = http.StatusUnsupportedMediaType, err.Error(), ErrValidation
	case errors.Is(err, NotImplementedError):
		code, message, errType = http.StatusNotImplemented, err.Error(), ErrNotImplemented
	case errors.Is(err, AuthenticationError):
		code, message, errType = http.StatusUnauthorized, err.Error(), ErrAuthentication
	case errors.Is(err, AuthorizationError):
		code, message, errType = http.StatusUnauthorized, err.Error(), ErrAuthorization
	case errors.Is(err, FatalError):
		code, message, errType = http.StatusInternalServerError, err.Error(), ErrFatal
	default:
		code, message, errType = http.StatusInternalServerError, err.Error(), ErrFatal
	}
	return code, message, errType
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
