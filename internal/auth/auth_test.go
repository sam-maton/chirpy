package auth

import (
	"net/http"
	"testing"
)

// type TestRun struct {
// 	params   map[string]any
// 	want     string
// 	hasError bool
// }

func TestGetBearerToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		header := http.Header{}
		header.Add("Authorization", "Bearer Success")

		got, err := GetBearerToken(header)

		if got != "Success" {
			t.Errorf("got %s, want Success", got)
		}

		if err != nil {
			t.Errorf("expected no error but got %s", err)
		}
	})
}
