package service

import (
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func TestValidateEndpointAllowsHTTPAndHTTPS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		wantCode string
	}{
		{name: "https origin", endpoint: "https://example.com", wantCode: ""},
		{name: "http origin", endpoint: "http://example.com", wantCode: ""},
		{name: "rejects ftp", endpoint: "ftp://example.com", wantCode: "CHANNEL_MONITOR_ENDPOINT_SCHEME"},
		{name: "rejects path", endpoint: "http://example.com/v1", wantCode: "CHANNEL_MONITOR_ENDPOINT_PATH"},
		{name: "rejects empty", endpoint: "", wantCode: "CHANNEL_MONITOR_INVALID_ENDPOINT"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateEndpoint(tt.endpoint)
			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("validateEndpoint(%q) = %v, want nil", tt.endpoint, err)
				}
				return
			}
			if got := infraerrors.Reason(err); got != tt.wantCode {
				t.Fatalf("validateEndpoint(%q) reason = %q, want %q (err=%v)", tt.endpoint, got, tt.wantCode, err)
			}
		})
	}
}
