package utils

import (
	"encoding/json"
	"log/slog"
)

// Converts a map[string]any to byte array
func ConvertMessageToJson(parsedMessage map[string]any) ([]byte, error) {
	jsonMessage, err := json.Marshal(parsedMessage)
	if err != nil {
		slog.Error("error: " + err.Error())
		return nil, err
	}
	return jsonMessage, nil
}

// Converts a byte array to a map[string]any
func ConvertJsonToMessage(jsonMessage []byte) (map[string]any, error) {
	var parsedMessage map[string]any
	err := json.Unmarshal(jsonMessage, &parsedMessage)
	if err != nil {
		slog.Error("error: " + err.Error())
		return nil, err
	}
	return parsedMessage, nil
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// Note: This is not cryptographically secure. For secure random, use crypto/rand.
		b[i] = charset[int64(i*37+length)%int64(len(charset))]
	}
	return string(b)
}
