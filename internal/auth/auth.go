package auth

import (
	"errors"
	"net/http"
	"strings"
)

func ParseApiKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")

	if header == "" {
		return "", errors.New("no api key provided")
	}

	apiKey := strings.TrimPrefix(header, "ApiKey ")

	return apiKey, nil
}
