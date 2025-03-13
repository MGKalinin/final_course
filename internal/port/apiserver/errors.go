package apiserver

// APIError struct обработка ошибок API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code int, message string) error {
	return APIError{Code: code, Message: message}
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) ToHTTPCode() int {
	return e.Code
}
