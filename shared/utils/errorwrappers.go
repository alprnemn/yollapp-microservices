package utils

import (
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", req.Method, req.URL.Path, err.Error())
	WriteError(w, http.StatusInternalServerError, err.Error())
}

func BadRequestResponse(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	WriteError(w, http.StatusBadRequest, err.Error())
}

func DatabaseError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("db error: %s path: %s error: %s", req.Method, req.URL.Path, err.Error())
	WriteError(w, http.StatusInternalServerError, err.Error())
}

func NotFoundError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	WriteError(w, http.StatusNotFound, err.Error())
}

func ConflictError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	WriteError(w, http.StatusConflict, err.Error())
}

func RateLimitExceededError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	log.Print("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	WriteError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	log.Print("unauthorized error ", "method ", r.Method, " path ", r.URL.Path, " error ", err.Error())
	WriteError(w, http.StatusUnauthorized, err.Error())
}
