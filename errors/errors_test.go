package errors

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type fieldError struct{}

func (f fieldError) Tag() string                       { return "" }
func (f fieldError) ActualTag() string                 { return "required" }
func (f fieldError) Namespace() string                 { return "" }
func (f fieldError) StructNamespace() string           { return "" }
func (f fieldError) Field() string                     { return "test" }
func (f fieldError) StructField() string               { return "" }
func (f fieldError) Param() string                     { return "test" }
func (f fieldError) Value() interface{}                { return "not test" }
func (f fieldError) Kind() reflect.Kind                { return reflect.String }
func (f fieldError) Type() reflect.Type                { return reflect.TypeOf("") }
func (f fieldError) Translate(ut ut.Translator) string { return "" }

func TestErrorHandler(t *testing.T) {
	type args struct {
		err error
	}
	type res struct {
		code    int
		message string
		errType ErrorType
	}

	tests := []struct {
		name string
		args args
		want res
	}{
		{
			name: "Return correct values for defined Validation error",
			args: args{
				err: ValidationError,
			},
			want: res{
				code:    http.StatusBadRequest,
				message: ValidationError.Error(),
				errType: ErrValidation,
			},
		},
		{
			name: "Return correct values for validator Validation error",
			args: args{
				err: validator.ValidationErrors{fieldError{}},
			},
			want: res{
				code:    http.StatusBadRequest,
				message: "Validation failed on field { test }, Condition: required{ test }, Value Recieved: not test",
				errType: ErrValidation,
			},
		},
		{
			name: "Return correct values for io EOF error",
			args: args{
				err: io.EOF,
			},
			want: res{
				code:    http.StatusBadRequest,
				message: "No request body",
				errType: ErrValidation,
			},
		},
		{
			name: "Return correct values for Not found error",
			args: args{
				err: NotFoundError,
			},
			want: res{
				code:    http.StatusNotFound,
				message: NotFoundError.Error(),
				errType: ErrNotFound,
			},
		},
		{
			name: "Return correct values for Entry exists error",
			args: args{
				err: EntryExistsError,
			},
			want: res{
				code:    http.StatusConflict,
				message: EntryExistsError.Error(),
				errType: ErrEntryExists,
			},
		},
		{
			name: "Return correct values for Entry deleted error",
			args: args{
				err: EntryDeletedError,
			},
			want: res{
				code:    http.StatusNotFound,
				message: EntryDeletedError.Error(),
				errType: ErrEntryDeleted,
			},
		},
		{
			name: "Return correct values for Invalid token error",
			args: args{
				err: InvalidTokenError,
			},
			want: res{
				code:    http.StatusUnauthorized,
				message: InvalidTokenError.Error(),
				errType: ErrInvalidToken,
			},
		},
		{
			name: "Return correct values for Expired token error",
			args: args{
				err: ExpiredTokenError,
			},
			want: res{
				code:    http.StatusUnauthorized,
				message: ExpiredTokenError.Error(),
				errType: ErrExpiredToken,
			},
		},
		{
			name: "Return correct values for Permission error",
			args: args{
				err: PermissionError,
			},
			want: res{
				code:    http.StatusForbidden,
				message: PermissionError.Error(),
				errType: ErrPermission,
			},
		},
		{
			name: "Return correct values for Unsupported media type error",
			args: args{
				err: UnsupportedMediaError,
			},
			want: res{
				code:    http.StatusUnsupportedMediaType,
				message: UnsupportedMediaError.Error(),
				errType: ErrValidation,
			},
		},
		{
			name: "Return correct values for Not implemented error",
			args: args{
				err: NotImplementedError,
			},
			want: res{
				code:    http.StatusNotImplemented,
				message: NotImplementedError.Error(),
				errType: ErrNotImplemented,
			},
		},
		{
			name: "Return correct values for Unauthorized error",
			args: args{
				err: AuthorizationError,
			},
			want: res{
				code:    http.StatusUnauthorized,
				message: AuthorizationError.Error(),
				errType: ErrAuthorization,
			},
		},
		{
			name: "Return correct values for Fatal error",
			args: args{
				err: FatalError,
			},
			want: res{
				code:    http.StatusInternalServerError,
				message: FatalError.Error(),
				errType: ErrFatal,
			},
		},
		{
			name: "Return correct values for bcrypt mismatch hash and password error",
			args: args{
				err: bcrypt.ErrMismatchedHashAndPassword,
			},
			want: res{
				code:    http.StatusUnauthorized,
				message: fmt.Sprintf("%s: Invalid login credentials", AuthenticationError.Error()),
				errType: ErrAuthentication,
			},
		},
		{
			name: "Return correct values for mongo write exception error",
			args: args{
				err: mongo.WriteException{},
			},
			want: res{
				code:    http.StatusConflict,
				message: EntryExistsError.Error(),
				errType: ErrEntryExists,
			},
		},
		{
			name: "Return correct value for wrapped error",
			args: args{
				err: fmt.Errorf("%w: Invalid login credentials", AuthenticationError),
			},
			want: res{
				code:    http.StatusUnauthorized,
				message: fmt.Sprintf("%s: Invalid login credentials", AuthenticationError.Error()),
				errType: ErrAuthentication,
			},
		},
		{
			name: "Return correct value for undefined error",
			args: args{
				fmt.Errorf("Undefined error"),
			},
			want: res{
				code:    http.StatusInternalServerError,
				message: "Undefined error",
				errType: ErrFatal,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, message, errType := ErrorHandler(test.args.err)
			Equal(t, test.want.code, code)
			Equal(t, test.want.message, message)
			Equal(t, test.want.errType, errType)
		})
	}
}
