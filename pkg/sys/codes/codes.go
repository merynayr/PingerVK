package codes

import "net/http"

// Code определяет тип кода ошибки
type Code int

// HTTP-коды ошибок, соответствующие стандартным кодам HTTP.
const (
	OK                  Code = http.StatusOK                  // 200 OK
	BadRequest          Code = http.StatusBadRequest          // 400 Bad Request
	Unauthorized        Code = http.StatusUnauthorized        // 401 Unauthorized
	Forbidden           Code = http.StatusForbidden           // 403 Forbidden
	NotFound            Code = http.StatusNotFound            // 404 Not Found
	MethodNotAllowed    Code = http.StatusMethodNotAllowed    // 405 Method Not Allowed
	Conflict            Code = http.StatusConflict            // 409 Conflict
	Gone                Code = http.StatusGone                // 410 Gone
	TooManyRequests     Code = http.StatusTooManyRequests     // 429 Too Many Requests
	InternalServerError Code = http.StatusInternalServerError // 500 Internal Server Error
	NotImplemented      Code = http.StatusNotImplemented      // 501 Not Implemented
	BadGateway          Code = http.StatusBadGateway          // 502 Bad Gateway
	ServiceUnavailable  Code = http.StatusServiceUnavailable  // 503 Service Unavailable
	GatewayTimeout      Code = http.StatusGatewayTimeout      // 504 Gateway Timeout
)
