package errors

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")

	ErrNotFound = errors.New("Your requested Item is not found")

	ErrUnAuthorize = errors.New("Unauthorize")

	ErrConflict = errors.New("Your Item already exist or duplicate")

	ErrBadParamInput = errors.New("Bad Request")

	ErrPublicKey = errors.New("invalid Public Key")

	ErrInvalidDataType = errors.New("invalid data type")

	ErrIsRequired = errors.New("is required")

	ErrInvalidValue = errors.New("invalid value")

	ErrForbidden = errors.New("you don't have permission to access this resource")
)

// authValidation
var (
	ErrInvalidToken             = errors.New("Invalid authorization token")
	ErrInvalidTokenType         = errors.New("Authorization token type does not match")
	ErrNotMatchTokenCredentials = errors.New("Authorization token credentials do not match")
	ErrInvalidTokenCredentials  = errors.New("Invalid authorization token credentials")
	ErrInvalidTokenExpired      = errors.New("Authorization token has expired")
)
