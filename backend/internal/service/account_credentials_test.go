package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountHasUsableCredentials(t *testing.T) {
	require.Error(t, accountHasUsableCredentials(&Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{},
	}))
	require.NoError(t, accountHasUsableCredentials(&Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"access_token": "tok",
		},
	}))
	require.NoError(t, accountHasUsableCredentials(&Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key": "sk",
		},
	}))
}

type credProbeRepo struct {
	AccountRepository
	account *Account
}

func (r *credProbeRepo) GetByID(ctx context.Context, id int64) (*Account, error) {
	return r.account, nil
}

type credProbeTester struct {
	result *ScheduledTestResult
}

func (t *credProbeTester) RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error) {
	return t.result, nil
}

func TestTestCredentialsFailClosedWithoutTester(t *testing.T) {
	svc := &AccountService{accountRepo: &credProbeRepo{account: &Account{
		ID:       1,
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"access_token": "x",
		},
	}}}
	err := svc.TestCredentials(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CREDENTIAL_TEST_UNAVAILABLE")
}

func TestTestCredentialsSurfacesProbeFailure(t *testing.T) {
	svc := &AccountService{accountRepo: &credProbeRepo{account: &Account{
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key": "sk-test",
		},
	}}}
	svc.SetConnectionTester(&credProbeTester{result: &ScheduledTestResult{
		Status:       "failed",
		ErrorMessage: "upstream 401",
	}})
	err := svc.TestCredentials(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "upstream 401")
}

func TestTestCredentialsSuccess(t *testing.T) {
	svc := &AccountService{accountRepo: &credProbeRepo{account: &Account{
		ID:       1,
		Platform: PlatformGemini,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"refresh_token": "r",
		},
	}}}
	svc.SetConnectionTester(&credProbeTester{result: &ScheduledTestResult{Status: "success"}})
	require.NoError(t, svc.TestCredentials(context.Background(), 1))
}
