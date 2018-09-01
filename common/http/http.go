package http

import (
	"net/http"
	"strings"
)

func Post(uri string, data string) int {
	req, _ := http.NewRequest(
		"POST",
		uri,
		strings.NewReader(data),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Status code 520 Unknown Error
		return 520
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
