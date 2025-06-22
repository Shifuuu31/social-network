package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JsonString(obj interface{}) (string, error) {
	var buffer bytes.Buffer

	enc := json.NewEncoder(&buffer)

	// Optional: pretty-print JSON
	enc.SetIndent("", "  ")

	if err := enc.Encode(obj); err != nil {
		return "", fmt.Errorf("JSON encode failed: %w", err)
	}

	return buffer.String(), nil
}

func JsonObj(jsonStr string) (obj interface{}, err error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))

	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&obj); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	// Optional: ensure there's nothing else in the JSON stream (e.g., extra commas)
	if decoder.More() {
		return nil, fmt.Errorf("unexpected data after JSON object")
	}

	return obj, nil
}
