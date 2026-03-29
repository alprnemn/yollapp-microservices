package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const MaxBytes = 1 << 20

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	er := map[string]string{
		"error": message,
	}
	_ = WriteJSON(w, status, er)
}

func ParseJSON(w http.ResponseWriter, req *http.Request, data any) error {

	req.Body = http.MaxBytesReader(w, req.Body, int64(MaxBytes))

	if req.Body == nil {
		return errors.New("missing request body")
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body must not be empty")
		}

		// Wrong type (e.g., string instead of number)
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			return fmt.Errorf("field '%s' must be of type %s", ute.Field, ute.Type.String())
		}

		// Syntax error
		var se *json.SyntaxError
		if errors.As(err, &se) {
			return fmt.Errorf("badly-formed JSON at position %d", se.Offset)
		}

		return err
	}

	// Ensure only one JSON object
	if decoder.More() {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}
