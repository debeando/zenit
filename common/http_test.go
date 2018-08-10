package common_test

import (
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/swapbyt3s/zenit/common"
)

func TestHTTPPost(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  }))
  defer ts.Close()

  response := common.HTTPPost(ts.URL, "foo")
  if ! response {
    t.Errorf("HTTPPost() return an error.")
  }
}
