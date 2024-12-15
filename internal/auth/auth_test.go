package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {

	tests := []struct {
		name      string
		want      string
		expectErr bool
		header    http.Header
	}{
		{
			name:      "Got correct response",
			want:      "Success",
			expectErr: false,
			header: http.Header{
				"Authorization": []string{"Bearer Success"},
			},
		},
		{
			name:      "Got incorrect response",
			want:      "",
			expectErr: true,
			header: http.Header{
				"Authorization": []string{"Incorrect_token"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.header)

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}

			if (err != nil) != tt.expectErr {
				t.Errorf("GetBearerToken() err = %v, expectErr = %v", err, tt.expectErr)
			}
		})
	}
}
