package errors

import "fmt"

// AppError adalah struktur dasar untuk error aplikasi
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// AuthError merepresentasikan error autentikasi
func AuthError(message string, err error) *AppError {
	return &AppError{
		Code:    401,
		Message: message,
		Err:     err,
	}
}

// BadRequestError merepresentasikan error permintaan yang tidak valid
func BadRequestError(message string, err error) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

// InternalServerError merepresentasikan error server internal
func InternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}
