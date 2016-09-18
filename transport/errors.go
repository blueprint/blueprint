package http

import "errors"

// Error values should consist of two lowercased words to keep API errors
// responses short and sweet.
const (
	AccountLocked       = "account locked"
	ExpiredSession      = "expired session"
	InvalidCreds        = "invalid credentials"
	InvalidEmail        = "invalid email"
	InvalidFile         = "invalid file"
	InvalidID           = "invalid id"
	InvalidInput        = "invalid input"
	InvalidJSON         = "invalid json"
	InvalidJSONP        = "invalid jsonp"
	InvalidXML          = "invalid xml"
	InvalidYAML         = "invalid yaml"
	InvalidParameter    = "invalid parameter"
	InvalidPassword     = "invalid password"
	InvalidPermission   = "invalid permission"
	InvalidRegistration = "invalid registration"
	InvalidSession      = "invalid session"
	InvalidUser         = "invalid user"
	MissingSession      = "missing session"
	MissingUser         = "missing user"
)

var (
	ErrAccountLocked       = errors.New(AccountLocked)
	ErrExpiredSession      = errors.New(ExpiredSession)
	ErrInvalidCreds        = errors.New(InvalidCreds)
	ErrInvalidEmail        = errors.New(InvalidEmail)
	ErrInvalidFile         = errors.New(InvalidFile)
	ErrInvalidID           = errors.New(InvalidID)
	ErrInvalidInput        = errors.New(InvalidInput)
	ErrInvalidJSON         = errors.New(InvalidJSON)
	ErrInvalidJSONP        = errors.New(InvalidJSONP)
	ErrInvalidXML          = errors.New(InvalidXML)
	ErrInvalidYAML         = errors.New(InvalidYAML)
	ErrInvalidParameter    = errors.New(InvalidParameter)
	ErrInvalidPassword     = errors.New(InvalidPassword)
	ErrInvalidPermission   = errors.New(InvalidPermission)
	ErrInvalidRegistration = errors.New(InvalidRegistration)
	ErrInvalidSession      = errors.New(InvalidSession)
	ErrInvalidUser         = errors.New(InvalidUser)
	ErrMissingSession      = errors.New(MissingSession)
	ErrMissingUser         = errors.New(MissingUser)
)

// ErrorMap is a map of error messages to HTTP status codes.
//
// These error messages are used by Response to match standard error messages
// with thier proper HTTP status codes.
var ErrorMap = map[error]int{
	ErrAccountLocked:       StatusForbidden,
	ErrExpiredSession:      StatusBadRequest,
	ErrInvalidCreds:        StatusBadRequest,
	ErrInvalidEmail:        StatusBadRequest,
	ErrInvalidFile:         StatusBadRequest,
	ErrInvalidID:           StatusBadRequest,
	ErrInvalidInput:        StatusBadRequest,
	ErrInvalidJSON:         StatusBadRequest,
	ErrInvalidJSONP:        StatusBadRequest,
	ErrInvalidXML:          StatusBadRequest,
	ErrInvalidYAML:         StatusBadRequest,
	ErrInvalidParameter:    StatusBadRequest,
	ErrInvalidPassword:     StatusBadRequest,
	ErrInvalidPermission:   StatusForbidden,
	ErrInvalidRegistration: StatusBadRequest,
	ErrInvalidSession:      StatusBadRequest,
	ErrInvalidUser:         StatusBadRequest,
	ErrMissingSession:      StatusBadRequest,
	ErrMissingUser:         StatusBadRequest,
}

// ApplyErrorCode inserts an error key and status code value into the ErrorMap.
func ApplyErrorCode(err error, statusCode int) bool {
	_, exists := ErrorMap[err]
	if exists {
		return false
	}
	ErrorMap[err] = statusCode
	return true
}
