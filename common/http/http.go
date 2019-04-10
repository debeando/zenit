package http

import (
	"log"
	"net/http"
	"strings"
)

func Post(uri string, data string, headers map[string]string) int {
	req, err := http.NewRequest(
		"POST",
		uri,
		strings.NewReader(data),
	)
	if err != nil {
		log.Printf("E! - %s\n", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("E! - %s\n", err)
		return 520 // Status code 520 Unknown Error
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
