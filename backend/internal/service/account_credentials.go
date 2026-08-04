package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// AccountConnectionTester probes whether an account's credentials can reach upstream.
// AccountTestService implements this via RunTestBackground.
type AccountConnectionTester interface {
	RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error)
}

// SetConnectionTester injects the upstream connection probe used by TestCredentials.
func (s *AccountService) SetConnectionTester(tester AccountConnectionTester) {
	if s == nil {
		return
	}
	s.connectionTester = tester
}

func accountHasUsableCredentials(account *Account) error {
	if account == nil {
		return ErrAccountNilInput
	}
	if len(account.Credentials) == 0 {
		return infraerrors.BadRequest("CREDENTIALS_EMPTY", "account credentials are empty")
	}

	switch account.Platform {
	case PlatformAnthropic, PlatformOpenAI, PlatformGemini, PlatformGrok, PlatformAntigravity:
		if account.Type == AccountTypeAPIKey || strings.EqualFold(account.Type, "apikey") {
			if strings.TrimSpace(account.GetCredential("api_key")) == "" &&
				strings.TrimSpace(account.GetCredential("access_token")) == "" {
				return infraerrors.BadRequest("CREDENTIALS_MISSING_API_KEY", "api_key or access_token is required")
			}
			return nil
		}
		if strings.TrimSpace(account.GetCredential("access_token")) == "" &&
			strings.TrimSpace(account.GetCredential("refresh_token")) == "" {
			return infraerrors.BadRequest("CREDENTIALS_MISSING_TOKEN", "access_token or refresh_token is required")
		}
		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", account.Platform)
	}
}

// TestCredentials validates stored credentials and probes upstream connectivity.
// It never returns success without a successful probe when a tester is configured.
func (s *AccountService) TestCredentials(ctx context.Context, id int64) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}
	if err := accountHasUsableCredentials(account); err != nil {
		return err
	}

	switch account.Platform {
	case PlatformAnthropic, PlatformOpenAI, PlatformGemini, PlatformGrok, PlatformAntigravity:
	default:
		return fmt.Errorf("unsupported platform: %s", account.Platform)
	}

	if s.connectionTester == nil {
		return infraerrors.New(http.StatusNotImplemented, "CREDENTIAL_TEST_UNAVAILABLE",
			"account connection tester is not configured; use admin account test endpoint")
	}

	result, err := s.connectionTester.RunTestBackground(ctx, id, "")
	if err != nil {
		return fmt.Errorf("credential probe failed: %w", err)
	}
	if result == nil {
		return infraerrors.BadRequest("CREDENTIAL_TEST_EMPTY_RESULT", "credential probe returned no result")
	}
	if result.Status != "success" {
		msg := strings.TrimSpace(result.ErrorMessage)
		if msg == "" {
			msg = "credential probe failed"
		}
		return infraerrors.BadRequest("CREDENTIAL_TEST_FAILED", msg)
	}
	return nil
}
