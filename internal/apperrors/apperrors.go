package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Internal      Type = "INTERNAL"
	BadRequest    Type = "BADREQUEST"
	NotFound      Type = "NOTFOUND"
	Authorization Type = "AUTHORIZATION"
)

// ServerError là message mặc định khi có lỗi hệ thống (DB down, bug...),
// không muốn lộ chi tiết thật (vd câu lỗi SQL) ra ngoài cho client.
const ServerError = "Something went wrong. Please try again later."

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case BadRequest:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case Authorization:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

// Parse luôn trả về *Error hợp lệ để handler chỉ cần 1 dòng là JSON được:
//
//	appErr := apperrors.Parse(err)
//	c.JSON(appErr.Status(), appErr)
//
// Nếu err không phải *Error (lỗi lạ từ DB/thư viện ngoài, chưa được wrap),
// ẩn chi tiết thật, trả về Internal chung chung.
func Parse(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return NewInternal()
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: ServerError,
	}
}

func NewNotFound(name, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}
