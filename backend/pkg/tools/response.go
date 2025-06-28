package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	const maxSize = 1048576 // 1MB
	limited := io.LimitReader(r.Body, maxSize)

	decoder := json.NewDecoder(limited) 
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return errors.New("invalid JSON: " + err.Error())
	}

	if decoder.More() {
		return errors.New("request must contain only a single JSON object")
	}

	return nil
}

func EncodeJSON(w http.ResponseWriter, status int, v any) (err error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(v); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(buffer.Bytes())

	return err
}

func RespondError(w http.ResponseWriter, msg string, status int) error {
	errResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	return EncodeJSON(w, status, errResp)
}
