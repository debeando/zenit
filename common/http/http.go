package http

import (
	"net/http"
	"strings"
)

func Post(uri string, data string, headers map[string]string) int {
	req, _ := http.NewRequest(
		"POST",
		uri,
		strings.NewReader(data),
	)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 520 // Status code 520 Unknown Error
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
