package http_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	test "github.com/swapbyt3s/zenit/common/http"
)

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		body, _ := ioutil.ReadAll(r.Body)

		if string(body) != "foo" {
			t.Errorf("Expected request to 'foo', got '%s'", string(body))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	response := test.Post(ts.URL, "foo")

	if response != 200 {
		t.Errorf("Expected: '200', got: '%v'.", response)
	}
}
