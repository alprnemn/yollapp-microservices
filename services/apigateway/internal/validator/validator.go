package validator

import (
	"github.com/alprnemn/yollapp-microservices/shared/utils"
	"net/http"
)

const MaxBytes = 15

func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost &&
			r.Method != http.MethodPut &&
			r.Method != http.MethodPatch {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			utils.WriteError(w, http.StatusBadRequest, "content type must be application/json")
			return
		}

		if r.Body == nil {
			utils.WriteError(w, http.StatusBadRequest, "missing request body")
			return
		}

		next.ServeHTTP(w, r)
	})
}
