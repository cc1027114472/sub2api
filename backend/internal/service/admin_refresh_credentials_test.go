package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefreshAccountCredentialsRejectsNonOAuth(t *testing.T) {
	svc := &adminServiceImpl{
		accountRepo: &credProbeRepo{account: &Account{
			ID:       1,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Credentials: map[string]any{
				"api_key": "sk",
			},
		}},
	}
	_, err := svc.RefreshAccountCredentials(t.Context(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "NOT_OAUTH")
}

func TestRefreshAccountCredentialsUnsupportedWithoutRefresher(t *testing.T) {
	svc := &adminServiceImpl{
		accountRepo: &credProbeRepo{account: &Account{
			ID:       1,
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"refresh_token": "r",
			},
		}},
	}
	_, err := svc.RefreshAccountCredentials(t.Context(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "REFRESH_UNSUPPORTED")
}
