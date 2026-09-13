//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestQuickSyncProbe_NormalizeModelID(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"models/gemini-2.5-flash", "gemini-2.5-flash"},
		{"models/gemini-2.5-pro", "gemini-2.5-pro"},
		{"gemini-2.5-flash", "gemini-2.5-flash"},
		{"claude-3-7-sonnet-20250219", "claude-3-7-sonnet-20250219"},
		{"  models/gpt-4o  ", "gpt-4o"},
		{"models/", ""},
		{"", ""},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			actual := NormalizeModelID(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQuickSyncProbe_NormalizeModelDisplayName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"models/gemini-2.5-flash", "gemini-2.5-flash"},
		{"Gemini 2.5 Pro", "Gemini 2.5 Pro"},
		{"  models/Claude 3.7 Sonnet  ", "Claude 3.7 Sonnet"},
		{"", ""},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			actual := NormalizeModelDisplayName(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQuickSyncProbe_AntigravityManager_V1Models(t *testing.T) {
	var authHeader, xAPIKeyHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models" {
			authHeader = r.Header.Get("Authorization")
			xAPIKeyHeader = r.Header.Get("x-api-key")

			resp := map[string]any{
				"object": "list",
				"data": []map[string]any{
					{"id": "models/gemini-2.5-flash", "display_name": "Gemini 2.5 Flash"},
					{"id": "gemini-2.5-pro", "display_name": "Gemini 2.5 Pro"},
					{"id": "claude-3-7-sonnet-20250219", "display_name": "Claude 3.7 Sonnet"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(server.Client())

	res, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL:  server.URL,
		APIKey:   "test-key-123",
		Platform: "antigravity",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 3, res.Total)
	require.Len(t, res.Models, 3)

	require.Equal(t, "Bearer test-key-123", authHeader)
	require.Equal(t, "test-key-123", xAPIKeyHeader)

	// Check model IDs are normalized
	modelMap := make(map[string]QuickSyncModelItem)
	for _, m := range res.Models {
		modelMap[m.ID] = m
		require.Equal(t, "token", m.BillingMode)
	}

	require.Contains(t, modelMap, "gemini-2.5-flash")
	require.Contains(t, modelMap, "gemini-2.5-pro")
	require.Contains(t, modelMap, "claude-3-7-sonnet-20250219")
	require.Equal(t, "Gemini 2.5 Flash", modelMap["gemini-2.5-flash"].DisplayName)
}

func TestQuickSyncProbe_AntigravityManager_V1BetaModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1beta/models" {
			resp := map[string]any{
				"models": []map[string]any{
					{
						"name":        "models/gemini-2.5-flash",
						"displayName": "models/Gemini 2.5 Flash",
					},
					{
						"name":        "models/gemini-2.5-pro",
						"displayName": "Gemini 2.5 Pro",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(server.Client())

	res, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL:  server.URL,
		APIKey:   "test-gemini-key",
		Platform: "gemini",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, res.Total)

	modelMap := make(map[string]QuickSyncModelItem)
	for _, m := range res.Models {
		modelMap[m.ID] = m
	}

	require.Contains(t, modelMap, "gemini-2.5-flash")
	require.Contains(t, modelMap, "gemini-2.5-pro")
	require.Equal(t, "Gemini 2.5 Flash", modelMap["gemini-2.5-flash"].DisplayName)
	require.Equal(t, "Gemini 2.5 Pro", modelMap["gemini-2.5-pro"].DisplayName)
}

func TestQuickSyncProbe_PricingLookup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data": []map[string]any{
				{"id": "gpt-5.4-mini"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	billingSvc := NewBillingService(nil, nil)
	svc := NewChannelQuickSyncService(nil, nil, billingSvc, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(server.Client())

	res, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Models, 1)

	model := res.Models[0]
	require.Equal(t, "gpt-5.4-mini", model.ID)
	require.NotNil(t, model.BasePriceIn)
	require.NotNil(t, model.BasePriceOut)
	require.Greater(t, *model.BasePriceIn, 0.0)
	require.Greater(t, *model.BasePriceOut, 0.0)
	require.Equal(t, *model.BasePriceIn, *model.PriceIn)
	require.Equal(t, *model.BasePriceOut, *model.PriceOut)
}

func TestQuickSyncProbe_ApplyBillingStrategy(t *testing.T) {
	inPrice := 0.000001
	outPrice := 0.000002

	items := []QuickSyncModelItem{
		{
			ID:           "test-model",
			DisplayName:  "Test Model",
			BasePriceIn:  &inPrice,
			BasePriceOut: &outPrice,
			BillingMode:  "token",
		},
	}

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	// 1. Test ratio multiplier
	ratioResult := svc.ApplyBillingStrategy(items, QuickSyncBillingStrategy{
		Mode:  "ratio",
		Ratio: 1.5,
	})
	require.Len(t, ratioResult, 1)
	require.Equal(t, "token", ratioResult[0].BillingMode)
	require.InDelta(t, 0.0000015, *ratioResult[0].PriceIn, 1e-10)
	require.InDelta(t, 0.000003, *ratioResult[0].PriceOut, 1e-10)
	require.Nil(t, ratioResult[0].PerReqPrice)

	// 2. Test per-request billing
	perReqPrice := 0.005
	perReqResult := svc.ApplyBillingStrategy(items, QuickSyncBillingStrategy{
		Mode:            "per_request",
		PerRequestPrice: &perReqPrice,
	})
	require.Len(t, perReqResult, 1)
	require.Equal(t, "per_request", perReqResult[0].BillingMode)
	require.NotNil(t, perReqResult[0].PerReqPrice)
	require.Equal(t, 0.005, *perReqResult[0].PerReqPrice)
	require.Nil(t, perReqResult[0].PriceIn)
	require.Nil(t, perReqResult[0].PriceOut)
}

func TestQuickSyncProbe_URLValidation(t *testing.T) {
	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	// Empty BaseURL
	_, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL: "",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "base_url is required")

	// Invalid URL scheme
	_, err = svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL: "ftp://invalid-scheme.com",
	})
	require.Error(t, err)
}

func TestQuickSyncProbe_FallbackToV1BetaOn404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path == "/v1beta/models" {
			resp := map[string]any{
				"models": []map[string]any{
					{"name": "models/gemini-2.5-flash"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(server.Client())

	res, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL:  server.URL,
		Platform: "antigravity", // platform default tries /v1/models then falls back to /v1beta/models
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 1, res.Total)
	require.Equal(t, "gemini-2.5-flash", res.Models[0].ID)
}

func TestQuickSyncProbe_Deduplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data": []map[string]any{
				{"id": "gemini-2.5-flash"},
				{"id": "models/gemini-2.5-flash"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(server.Client())

	res, err := svc.ProbeUpstream(context.Background(), QuickSyncProbeParams{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 1, res.Total)
	require.Equal(t, "gemini-2.5-flash", res.Models[0].ID)
}

// --- QuickSync Commit Tests & Mocks ---

type mockQuickSyncAccountRepo struct {
	AccountRepository
	createFn     func(ctx context.Context, account *Account) error
	bindGroupsFn func(ctx context.Context, accountID int64, groupIDs []int64) error
	created      *Account
	boundGroups  map[int64][]int64
}

func (m *mockQuickSyncAccountRepo) Create(ctx context.Context, account *Account) error {
	if m.createFn != nil {
		return m.createFn(ctx, account)
	}
	account.ID = 1001
	m.created = account
	return nil
}

func (m *mockQuickSyncAccountRepo) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	if m.bindGroupsFn != nil {
		return m.bindGroupsFn(ctx, accountID, groupIDs)
	}
	if m.boundGroups == nil {
		m.boundGroups = make(map[int64][]int64)
	}
	m.boundGroups[accountID] = groupIDs
	return nil
}

func (m *mockQuickSyncAccountRepo) GetByID(ctx context.Context, id int64) (*Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ExistsByID(ctx context.Context, id int64) (bool, error) { return true, nil }
func (m *mockQuickSyncAccountRepo) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) FindByExtraField(ctx context.Context, key string, value any) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListCRSAccountIDs(ctx context.Context) (map[string]int64, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) Update(ctx context.Context, account *Account) error { return nil }
func (m *mockQuickSyncAccountRepo) Delete(ctx context.Context, id int64) error { return nil }
func (m *mockQuickSyncAccountRepo) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) { return nil, nil, nil }
func (m *mockQuickSyncAccountRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) { return nil, nil, nil }
func (m *mockQuickSyncAccountRepo) ListAllWithFilters(ctx context.Context, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListActive(ctx context.Context) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListByPlatform(ctx context.Context, platform string) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) UpdateLastUsed(ctx context.Context, id int64) error { return nil }
func (m *mockQuickSyncAccountRepo) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error { return nil }
func (m *mockQuickSyncAccountRepo) SetError(ctx context.Context, id int64, errorMsg string) error { return nil }
func (m *mockQuickSyncAccountRepo) ClearError(ctx context.Context, id int64) error { return nil }
func (m *mockQuickSyncAccountRepo) SetSchedulable(ctx context.Context, id int64, schedulable bool) error { return nil }
func (m *mockQuickSyncAccountRepo) AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error) { return 0, nil }
func (m *mockQuickSyncAccountRepo) ListSchedulable(ctx context.Context) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) { return nil, nil }
func (m *mockQuickSyncAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) { return nil, nil }

type mockQuickSyncGroupRepo struct {
	GroupRepository
	groups         map[int64]*Group
	nextID         int64
	createFn       func(ctx context.Context, group *Group) error
	getByIDFn      func(ctx context.Context, id int64) (*Group, error)
	updateFn       func(ctx context.Context, group *Group) error
	existsByNameFn func(ctx context.Context, name string) (bool, error)
}

func newMockQuickSyncGroupRepo() *mockQuickSyncGroupRepo {
	return &mockQuickSyncGroupRepo{
		groups: make(map[int64]*Group),
		nextID: 50,
	}
}

func (m *mockQuickSyncGroupRepo) Create(ctx context.Context, group *Group) error {
	if m.createFn != nil {
		return m.createFn(ctx, group)
	}
	m.nextID++
	group.ID = m.nextID
	m.groups[group.ID] = group
	return nil
}

func (m *mockQuickSyncGroupRepo) GetByID(ctx context.Context, id int64) (*Group, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	g, ok := m.groups[id]
	if !ok {
		return nil, ErrGroupNotFound
	}
	return g, nil
}

func (m *mockQuickSyncGroupRepo) GetByIDLite(ctx context.Context, id int64) (*Group, error) {
	return m.GetByID(ctx, id)
}

func (m *mockQuickSyncGroupRepo) Update(ctx context.Context, group *Group) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, group)
	}
	m.groups[group.ID] = group
	return nil
}

func (m *mockQuickSyncGroupRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.existsByNameFn != nil {
		return m.existsByNameFn(ctx, name)
	}
	for _, g := range m.groups {
		if g.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockQuickSyncGroupRepo) Delete(ctx context.Context, id int64) error { return nil }
func (m *mockQuickSyncGroupRepo) DeleteCascade(ctx context.Context, id int64) ([]int64, error) { return nil, nil }
func (m *mockQuickSyncGroupRepo) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) { return nil, nil, nil }
func (m *mockQuickSyncGroupRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) { return nil, nil, nil }
func (m *mockQuickSyncGroupRepo) ListActive(ctx context.Context) ([]Group, error) { return nil, nil }
func (m *mockQuickSyncGroupRepo) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) { return nil, nil }
func (m *mockQuickSyncGroupRepo) GetAccountCount(ctx context.Context, groupID int64) (int64, int64, error) { return 0, 0, nil }
func (m *mockQuickSyncGroupRepo) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) { return 0, nil }
func (m *mockQuickSyncGroupRepo) GetAccountIDsByGroupIDs(ctx context.Context, groupIDs []int64) ([]int64, error) { return nil, nil }
func (m *mockQuickSyncGroupRepo) BindAccountsToGroup(ctx context.Context, groupID int64, accountIDs []int64) error { return nil }
func (m *mockQuickSyncGroupRepo) UpdateSortOrders(ctx context.Context, updates []GroupSortOrderUpdate) error { return nil }

type mockQuickSyncCacheInvalidator struct {
	invalidated bool
}

func (m *mockQuickSyncCacheInvalidator) InvalidateCache() {
	m.invalidated = true
}

func TestQuickSyncCommit_SuccessWithExistingGroups(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockChannelRepository{}
	cacheInvalidator := &mockQuickSyncCacheInvalidator{}

	// Setup existing groups 10 and 20
	groupRepo.groups[10] = &Group{
		ID:             10,
		Name:           "Group 10",
		Platform:       PlatformAntigravity,
		ModelAllowlist: GroupModelAllowlist{Enabled: true, Models: []string{}},
	}
	groupRepo.groups[20] = &Group{
		ID:             20,
		Name:           "Group 20",
		Platform:       PlatformAntigravity,
		ModelAllowlist: GroupModelAllowlist{Enabled: true, Models: []string{"claude-3-opus"}},
	}

	var createdChannel *Channel
	channelRepo.createFn = func(ctx context.Context, ch *Channel) error {
		ch.ID = 501
		createdChannel = ch
		return nil
	}
	channelRepo.existsByNameFn = func(ctx context.Context, name string) (bool, error) {
		return false, nil
	}
	channelRepo.getGroupsInOtherChannelsFn = func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
		return nil, nil
	}

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, cacheInvalidator)

	inPrice := 0.075
	outPrice := 0.30
	perReqPrice := 0.01

	params := QuickSyncCommitParams{
		Name:     "Test AG Channel",
		BaseURL:  "https://api.example.com",
		APIKey:   "sk-test-secret",
		Platform: "antigravity",
		Models: []QuickSyncCommitModelItem{
			{
				Model:         "gemini-2.5-flash",
				TargetGroupID: 10,
				BillingMode:   "token",
				InputPrice:    &inPrice,
				OutputPrice:   &outPrice,
			},
			{
				Model:           "claude-3-5-sonnet",
				TargetGroupID:   20,
				BillingMode:     "per_request",
				PerRequestPrice: &perReqPrice,
			},
		},
	}

	res, err := svc.CommitQuickSync(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, int64(501), res.ChannelID)
	require.Equal(t, int64(1001), res.AccountID)
	require.Equal(t, []int64{10, 20}, res.GroupIDs)
	require.Equal(t, 2, res.ModelCount)

	// Verify Account
	require.NotNil(t, accountRepo.created)
	require.Equal(t, "Test AG Channel", accountRepo.created.Name)
	require.Equal(t, "antigravity", accountRepo.created.Platform)
	require.Equal(t, AccountTypeAPIKey, accountRepo.created.Type)
	require.Equal(t, "sk-test-secret", accountRepo.created.Credentials["api_key"])
	require.Equal(t, "https://api.example.com", accountRepo.created.Credentials["base_url"])
	require.Equal(t, []int64{10, 20}, accountRepo.boundGroups[accountRepo.created.ID])

	// Verify Group ModelAllowlist updates
	require.Contains(t, groupRepo.groups[10].ModelAllowlist.Models, "gemini-2.5-flash")
	require.Contains(t, groupRepo.groups[20].ModelAllowlist.Models, "claude-3-opus")
	require.Contains(t, groupRepo.groups[20].ModelAllowlist.Models, "claude-3-5-sonnet")

	// Verify Channel & Pricing
	require.NotNil(t, createdChannel)
	require.Equal(t, "Test AG Channel", createdChannel.Name)
	require.Equal(t, []int64{10, 20}, createdChannel.GroupIDs)
	require.Len(t, createdChannel.ModelPricing, 2)

	p1 := createdChannel.ModelPricing[0]
	require.Equal(t, []string{"gemini-2.5-flash"}, p1.Models)
	require.Equal(t, BillingModeToken, p1.BillingMode)
	require.Equal(t, &inPrice, p1.InputPrice)
	require.Equal(t, &outPrice, p1.OutputPrice)

	p2 := createdChannel.ModelPricing[1]
	require.Equal(t, []string{"claude-3-5-sonnet"}, p2.Models)
	require.Equal(t, BillingModePerRequest, p2.BillingMode)
	require.Equal(t, 0.0, *p2.InputPrice)
	require.Equal(t, 0.0, *p2.OutputPrice)
	require.Equal(t, &perReqPrice, p2.PerRequestPrice)

	// Verify cache invalidation
	require.True(t, cacheInvalidator.invalidated)
}

func TestQuickSyncCommit_SuccessWithNewGroup(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockChannelRepository{}
	cacheInvalidator := &mockQuickSyncCacheInvalidator{}

	channelRepo.createFn = func(ctx context.Context, ch *Channel) error {
		ch.ID = 701
		return nil
	}
	channelRepo.existsByNameFn = func(ctx context.Context, name string) (bool, error) {
		return false, nil
	}
	channelRepo.getGroupsInOtherChannelsFn = func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
		return nil, nil
	}

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, cacheInvalidator)

	inPrice := 0.05
	outPrice := 0.15

	params := QuickSyncCommitParams{
		Name:     "Brand New Channel",
		BaseURL:  "https://api.example.com",
		APIKey:   "sk-brand-new",
		Platform: "antigravity",
		NewGroup: &QuickSyncNewGroupParams{
			Create:         true,
			Name:           "Custom Auto Group",
			RateMultiplier: 1.25,
		},
		Models: []QuickSyncCommitModelItem{
			{
				Model:         "gpt-4o",
				TargetGroupID: 0, // Should resolve to new group
				BillingMode:   "token",
				InputPrice:    &inPrice,
				OutputPrice:   &outPrice,
			},
		},
	}

	res, err := svc.CommitQuickSync(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, int64(701), res.ChannelID)
	require.Equal(t, int64(1001), res.AccountID)
	require.Len(t, res.GroupIDs, 1)
	newGID := res.GroupIDs[0]
	require.True(t, newGID > 0)
	require.Equal(t, 1, res.ModelCount)

	// Verify newly created group
	createdGroup, ok := groupRepo.groups[newGID]
	require.True(t, ok)
	require.Equal(t, "Custom Auto Group", createdGroup.Name)
	require.Equal(t, 1.25, createdGroup.RateMultiplier)
	require.Equal(t, "antigravity", createdGroup.Platform)
	require.Contains(t, createdGroup.ModelAllowlist.Models, "gpt-4o")

	// Verify Account bound to new group
	require.Equal(t, []int64{newGID}, accountRepo.boundGroups[accountRepo.created.ID])
}

func TestQuickSyncCommit_ValidationErrors(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockChannelRepository{}
	svc := NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, nil)

	// 1. Missing Name
	_, err := svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 1}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "name is required")

	// 2. Missing BaseURL
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 1}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "base_url is required")

	// 3. Missing APIKey
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 1}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "api_key is required")

	// 4. Empty Models
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  nil,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "models cannot be empty")

	// 5. Target group not found
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 9999}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "get group 9999")

	// 6. Target group requires OAuth only
	groupRepo.groups[88] = &Group{
		ID:               88,
		Name:             "OAuth Only Group",
		Platform:         PlatformAntigravity,
		RequireOAuthOnly: true,
	}
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 88}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "requires OAuth accounts")

	// 7. Channel name already exists
	groupRepo.groups[10] = &Group{ID: 10, Name: "Group 10", Platform: PlatformAntigravity}
	channelRepo.existsByNameFn = func(ctx context.Context, name string) (bool, error) {
		return true, nil
	}
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Existing Channel",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 10}},
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrChannelExists))
	channelRepo.existsByNameFn = func(ctx context.Context, name string) (bool, error) { return false, nil }

	// 8. Group already in another channel
	channelRepo.getGroupsInOtherChannelsFn = func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
		return []int64{10}, nil
	}
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models:  []QuickSyncCommitModelItem{{Model: "gpt-4o", TargetGroupID: 10}},
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrGroupAlreadyInChannel))
	channelRepo.getGroupsInOtherChannelsFn = func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) { return nil, nil }

	// 9. Negative price
	negPrice := -0.5
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models: []QuickSyncCommitModelItem{
			{Model: "gpt-4o", TargetGroupID: 10, InputPrice: &negPrice},
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "input_price must be >= 0")

	// 10. Per request mode missing per_request_price
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models: []QuickSyncCommitModelItem{
			{Model: "gpt-4o", TargetGroupID: 10, BillingMode: "per_request"},
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "per_request_price required")

	// 11. Duplicate models
	_, err = svc.CommitQuickSync(context.Background(), QuickSyncCommitParams{
		Name:    "Chan",
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Models: []QuickSyncCommitModelItem{
			{Model: "gpt-4o", TargetGroupID: 10},
			{Model: "models/gpt-4o", TargetGroupID: 10},
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate model")
}

func TestQuickSyncCommit_ModelPricingConstruction(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockChannelRepository{}

	groupRepo.groups[1] = &Group{ID: 1, Name: "G1", Platform: PlatformAntigravity}
	var capturedPricing []ChannelModelPricing
	channelRepo.createFn = func(ctx context.Context, ch *Channel) error {
		capturedPricing = ch.ModelPricing
		return nil
	}
	channelRepo.existsByNameFn = func(ctx context.Context, name string) (bool, error) { return false, nil }
	channelRepo.getGroupsInOtherChannelsFn = func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) { return nil, nil }

	svc := NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, nil)

	inPrice := 0.01
	outPrice := 0.02
	perReqPrice := 0.005

	params := QuickSyncCommitParams{
		Name:     "Pricing Test Channel",
		BaseURL:  "https://api.example.com",
		APIKey:   "sk-key",
		Platform: "antigravity",
		Models: []QuickSyncCommitModelItem{
			{
				Model:         "token-model",
				TargetGroupID: 1,
				BillingMode:   "token",
				InputPrice:    &inPrice,
				OutputPrice:   &outPrice,
			},
			{
				Model:           "per-req-model",
				TargetGroupID:   1,
				BillingMode:     "per_request",
				PerRequestPrice: &perReqPrice,
			},
		},
	}

	_, err := svc.CommitQuickSync(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, capturedPricing, 2)

	// Item 1: Token mode
	require.Equal(t, BillingModeToken, capturedPricing[0].BillingMode)
	require.Equal(t, &inPrice, capturedPricing[0].InputPrice)
	require.Equal(t, &outPrice, capturedPricing[0].OutputPrice)
	require.Nil(t, capturedPricing[0].PerRequestPrice)

	// Item 2: PerRequest mode
	require.Equal(t, BillingModePerRequest, capturedPricing[1].BillingMode)
	require.Equal(t, 0.0, *capturedPricing[1].InputPrice)
	require.Equal(t, 0.0, *capturedPricing[1].OutputPrice)
	require.Equal(t, &perReqPrice, capturedPricing[1].PerRequestPrice)
}


