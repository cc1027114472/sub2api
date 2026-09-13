//go:build unit

package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// mockQuickSyncAccountRepo implements minimal service.AccountRepository for handler testing.
type mockQuickSyncAccountRepo struct {
	service.AccountRepository
	createdAccount *service.Account
	boundGroups    map[int64][]int64
}

func (m *mockQuickSyncAccountRepo) Create(ctx context.Context, account *service.Account) error {
	account.ID = 1001
	m.createdAccount = account
	return nil
}

func (m *mockQuickSyncAccountRepo) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	if m.boundGroups == nil {
		m.boundGroups = make(map[int64][]int64)
	}
	m.boundGroups[accountID] = groupIDs
	return nil
}

// mockQuickSyncGroupRepo implements minimal service.GroupRepository for handler testing.
type mockQuickSyncGroupRepo struct {
	service.GroupRepository
	groups map[int64]*service.Group
}

func newMockQuickSyncGroupRepo() *mockQuickSyncGroupRepo {
	return &mockQuickSyncGroupRepo{
		groups: make(map[int64]*service.Group),
	}
}

func (m *mockQuickSyncGroupRepo) GetByID(ctx context.Context, id int64) (*service.Group, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, service.ErrGroupNotFound
	}
	return g, nil
}

func (m *mockQuickSyncGroupRepo) GetByIDLite(ctx context.Context, id int64) (*service.Group, error) {
	return m.GetByID(ctx, id)
}

func (m *mockQuickSyncGroupRepo) Create(ctx context.Context, group *service.Group) error {
	group.ID = 2001
	m.groups[group.ID] = group
	return nil
}

func (m *mockQuickSyncGroupRepo) Update(ctx context.Context, group *service.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockQuickSyncGroupRepo) UpdateModelAllowlist(ctx context.Context, groupID int64, allowlist service.GroupModelAllowlist) error {
	if g, ok := m.groups[groupID]; ok {
		g.ModelAllowlist = allowlist
	}
	return nil
}

// mockQuickSyncChannelRepo implements minimal service.ChannelRepository for handler testing.
type mockQuickSyncChannelRepo struct {
	service.ChannelRepository
	createdChannel *service.Channel
}

func (m *mockQuickSyncChannelRepo) Create(ctx context.Context, ch *service.Channel) error {
	ch.ID = 3001
	m.createdChannel = ch
	return nil
}

func (m *mockQuickSyncChannelRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (m *mockQuickSyncChannelRepo) GetGroupsInOtherChannels(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
	return nil, nil
}

func setupQuickSyncTestRouter(svc *service.ChannelQuickSyncService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewChannelHandler(nil, nil, nil)
	if svc != nil {
		handler.SetQuickSyncService(svc)
	}

	router := gin.New()
	channels := router.Group("/api/v1/admin/channels")
	{
		channels.POST("/quick-sync/probe", handler.QuickSyncProbe)
		channels.POST("/quick-sync/commit", handler.QuickSyncCommit)
	}
	return router
}

func TestQuickSyncProbe_Success(t *testing.T) {
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"models": []map[string]any{
				{
					"id":           "models/gemini-2.5-flash",
					"display_name": "Gemini 2.5 Flash",
				},
				{
					"id":           "models/gemini-2.5-pro",
					"display_name": "Gemini 2.5 Pro",
				},
			},
		})
	}))
	defer mockUpstream.Close()

	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	svc.SetHTTPClient(mockUpstream.Client())

	router := setupQuickSyncTestRouter(svc)

	payload := map[string]any{
		"base_url": mockUpstream.URL,
		"api_key":  "test-secret-key",
		"platform": "antigravity",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/probe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code int                         `json:"code"`
		Data service.QuickSyncProbeResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, 2, resp.Data.Total)
	require.Len(t, resp.Data.Models, 2)
	require.Equal(t, "gemini-2.5-flash", resp.Data.Models[0].ID)
	require.Equal(t, "gemini-2.5-pro", resp.Data.Models[1].ID)
}

func TestQuickSyncProbe_MissingBaseURL(t *testing.T) {
	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	router := setupQuickSyncTestRouter(svc)

	payload := map[string]any{
		"api_key":  "test-secret-key",
		"platform": "antigravity",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/probe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuickSyncProbe_InvalidJSON(t *testing.T) {
	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	router := setupQuickSyncTestRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/probe", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuickSyncProbe_ServiceUnavailable(t *testing.T) {
	router := setupQuickSyncTestRouter(nil)

	payload := map[string]any{
		"base_url": "https://api.example.com",
		"api_key":  "test-secret-key",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/probe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQuickSyncCommit_Success(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockQuickSyncChannelRepo{}

	groupRepo.groups[10] = &service.Group{
		ID:             10,
		Name:           "Default Group",
		Platform:       service.PlatformAntigravity,
		ModelAllowlist: service.GroupModelAllowlist{Enabled: true, Models: []string{}},
	}

	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, nil)
	router := setupQuickSyncTestRouter(svc)

	inPrice := 0.075
	outPrice := 0.30

	commitPayload := map[string]any{
		"name":     "Antigravity Node 1",
		"base_url": "https://api.example.com",
		"api_key":  "sk-test-secret",
		"platform": "antigravity",
		"models": []map[string]any{
			{
				"model":           "gemini-2.5-flash",
				"target_group_id": 10,
				"billing_mode":   "token",
				"input_price":     inPrice,
				"output_price":    outPrice,
			},
		},
	}
	body, err := json.Marshal(commitPayload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/commit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code int                          `json:"code"`
		Data service.QuickSyncCommitResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, int64(3001), resp.Data.ChannelID)
	require.Equal(t, int64(1001), resp.Data.AccountID)
	require.Equal(t, 1, resp.Data.ModelCount)
	require.Contains(t, resp.Data.GroupIDs, int64(10))
}

func TestQuickSyncCommit_ValidationErrors(t *testing.T) {
	accountRepo := &mockQuickSyncAccountRepo{}
	groupRepo := newMockQuickSyncGroupRepo()
	channelRepo := &mockQuickSyncChannelRepo{}
	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, accountRepo, groupRepo, channelRepo, nil, nil)
	router := setupQuickSyncTestRouter(svc)

	tests := []struct {
		name    string
		payload map[string]any
	}{
		{
			name: "missing name",
			payload: map[string]any{
				"base_url": "https://api.example.com",
				"api_key":  "sk-test",
				"models":   []map[string]any{{"model": "gpt-4"}},
			},
		},
		{
			name: "missing base_url",
			payload: map[string]any{
				"name":    "Channel 1",
				"api_key": "sk-test",
				"models":  []map[string]any{{"model": "gpt-4"}},
			},
		},
		{
			name: "missing api_key",
			payload: map[string]any{
				"name":     "Channel 1",
				"base_url": "https://api.example.com",
				"models":   []map[string]any{{"model": "gpt-4"}},
			},
		},
		{
			name: "empty models",
			payload: map[string]any{
				"name":     "Channel 1",
				"base_url": "https://api.example.com",
				"api_key":  "sk-test",
				"models":   []map[string]any{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/commit", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestQuickSyncCommit_InvalidJSON(t *testing.T) {
	svc := service.NewChannelQuickSyncService(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	router := setupQuickSyncTestRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/commit", bytes.NewBufferString("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuickSyncCommit_ServiceUnavailable(t *testing.T) {
	router := setupQuickSyncTestRouter(nil)

	payload := map[string]any{
		"name":     "Node 1",
		"base_url": "https://api.example.com",
		"api_key":  "sk-test",
		"models":   []map[string]any{{"model": "gpt-4"}},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/channels/quick-sync/commit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
