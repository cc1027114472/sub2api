package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
)

// QuickSyncProbeParams holds parameters for probing upstream models.
type QuickSyncProbeParams struct {
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Platform string `json:"platform"` // default "antigravity"
}

// QuickSyncModelItem represents a discovered model with base and calculated pricing.
type QuickSyncModelItem struct {
	ID            string   `json:"id"`
	DisplayName   string   `json:"display_name"`
	BasePriceIn   *float64 `json:"base_price_in"`
	BasePriceOut  *float64 `json:"base_price_out"`
	BillingMode   string   `json:"billing_mode"` // "token" | "per_request"
	TargetGroupID int64    `json:"target_group_id"`
	PriceIn       *float64 `json:"price_in"`
	PriceOut      *float64 `json:"price_out"`
	PerReqPrice   *float64 `json:"per_request_price"`
}

// QuickSyncProbeResult is returned after a successful upstream model probe.
type QuickSyncProbeResult struct {
	Models   []QuickSyncModelItem `json:"models"`
	Total    int                  `json:"total"`
	Warnings []string             `json:"warnings,omitempty"`
}

// QuickSyncBillingStrategy configures global or per-model price calculations.
type QuickSyncBillingStrategy struct {
	Mode            string   `json:"mode"` // "ratio" | "per_request" | "fixed"
	Ratio           float64  `json:"ratio"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

// QuickSyncNewGroupParams defines parameters for creating a new group during quick sync.
type QuickSyncNewGroupParams struct {
	Create         bool    `json:"create"`
	Name           string  `json:"name"`
	RateMultiplier float64 `json:"rate_multiplier"`
}

// QuickSyncCommitModelItem defines per-model target group and pricing settings.
type QuickSyncCommitModelItem struct {
	Model           string   `json:"model"`
	TargetGroupID   int64    `json:"target_group_id"`
	BillingMode     string   `json:"billing_mode"` // "token" | "per_request"
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

// QuickSyncCommitParams defines the payload for committing a quick sync onboarding.
type QuickSyncCommitParams struct {
	Name            string                     `json:"name"`
	BaseURL         string                     `json:"base_url"`
	APIKey          string                     `json:"api_key"`
	Platform        string                     `json:"platform"`
	DefaultGroupID  *int64                     `json:"default_group_id"`
	NewGroup        *QuickSyncNewGroupParams   `json:"new_group"`
	BillingStrategy *QuickSyncBillingStrategy  `json:"billing_strategy"`
	Models          []QuickSyncCommitModelItem `json:"models"`
}

// QuickSyncCommitResult is returned after a successful quick sync commit.
type QuickSyncCommitResult struct {
	ChannelID  int64   `json:"channel_id"`
	AccountID  int64   `json:"account_id"`
	GroupIDs   []int64 `json:"group_ids"`
	ModelCount int     `json:"model_count"`
}

// ChannelQuickSyncService handles probing upstream models and quick sync onboarding.
type ChannelQuickSyncService struct {
	cfg                     *config.Config
	httpClient              *http.Client
	pricingService          *PricingService
	billingService          *BillingService
	accountTestService      *AccountTestService
	accountRepo             AccountRepository
	groupRepo               GroupRepository
	channelRepo             ChannelRepository
	channelService          *ChannelService
	entClient               *dbent.Client
	channelCacheInvalidator ChannelCacheInvalidator
	cachePubSub             ChannelCachePubSub
}

// NewChannelQuickSyncService creates a new ChannelQuickSyncService instance.
func NewChannelQuickSyncService(
	cfg *config.Config,
	pricingService *PricingService,
	billingService *BillingService,
	accountTestService *AccountTestService,
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	channelRepo ChannelRepository,
	entClient *dbent.Client,
	channelCacheInvalidator ChannelCacheInvalidator,
) *ChannelQuickSyncService {
	return &ChannelQuickSyncService{
		cfg:                     cfg,
		pricingService:          pricingService,
		billingService:          billingService,
		accountTestService:      accountTestService,
		accountRepo:             accountRepo,
		groupRepo:               groupRepo,
		channelRepo:             channelRepo,
		entClient:               entClient,
		channelCacheInvalidator: channelCacheInvalidator,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SetHTTPClient allows overriding the HTTP client (e.g. for testing).
func (s *ChannelQuickSyncService) SetHTTPClient(client *http.Client) {
	s.httpClient = client
}

// SetRepositories sets the repository dependencies for quick sync commit.
func (s *ChannelQuickSyncService) SetRepositories(accountRepo AccountRepository, groupRepo GroupRepository, channelRepo ChannelRepository) {
	s.accountRepo = accountRepo
	s.groupRepo = groupRepo
	s.channelRepo = channelRepo
}

// SetEntClient sets the ent client used for atomic transactions.
func (s *ChannelQuickSyncService) SetEntClient(client *dbent.Client) {
	s.entClient = client
}

// SetChannelCacheInvalidator sets the cache invalidator.
func (s *ChannelQuickSyncService) SetChannelCacheInvalidator(invalidator ChannelCacheInvalidator) {
	s.channelCacheInvalidator = invalidator
}

// SetCachePubSub sets the channel cache pub/sub.
func (s *ChannelQuickSyncService) SetCachePubSub(cachePubSub ChannelCachePubSub) {
	s.cachePubSub = cachePubSub
}

// SetChannelService sets the channel service.
func (s *ChannelQuickSyncService) SetChannelService(channelService *ChannelService) {
	s.channelService = channelService
}

// NormalizeModelID strips whitespace and upstream prefixes such as "models/".
func NormalizeModelID(raw string) string {
	id := strings.TrimSpace(raw)
	id = strings.TrimPrefix(id, "models/")
	return strings.TrimSpace(id)
}

// NormalizeModelDisplayName normalizes display names by removing leading "models/" prefix.
func NormalizeModelDisplayName(raw string) string {
	name := strings.TrimSpace(raw)
	name = strings.TrimPrefix(name, "models/")
	return strings.TrimSpace(name)
}

// validateUpstreamBaseURL performs SSRF and URL format validation against security config.
func (s *ChannelQuickSyncService) validateUpstreamBaseURL(raw string) (string, error) {
	if s.cfg == nil {
		return urlvalidator.ValidateURLFormat(raw, true)
	}
	if !s.cfg.Security.URLAllowlist.Enabled {
		return urlvalidator.ValidateURLFormat(raw, s.cfg.Security.URLAllowlist.AllowInsecureHTTP)
	}
	normalized, err := urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowedHosts:     s.cfg.Security.URLAllowlist.UpstreamHosts,
		RequireAllowlist: true,
		AllowPrivate:     s.cfg.Security.URLAllowlist.AllowPrivateHosts,
	})
	if err != nil {
		return "", err
	}
	return normalized, nil
}

type quickSyncRawModelEntry struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	DisplayNameSnake string `json:"display_name"`
	DisplayNameCamel string `json:"displayName"`
}

// ProbeUpstream requests upstream models and pairs them with official reference pricing.
func (s *ChannelQuickSyncService) ProbeUpstream(ctx context.Context, params QuickSyncProbeParams) (*QuickSyncProbeResult, error) {
	baseURL := strings.TrimSpace(params.BaseURL)
	if baseURL == "" {
		return nil, errors.New("base_url is required")
	}

	platform := strings.ToLower(strings.TrimSpace(params.Platform))
	if platform == "" {
		platform = PlatformAntigravity
	}

	normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base_url: %w", err)
	}

	apiKey := strings.TrimSpace(params.APIKey)

	// Determine candidate URLs to probe based on platform
	var candidateURLs []string
	switch platform {
	case PlatformGemini:
		candidateURLs = []string{
			buildGeminiModelsURL(normalizedBaseURL),
			buildV1ModelsURL(normalizedBaseURL),
		}
	default:
		candidateURLs = []string{
			buildV1ModelsURL(normalizedBaseURL),
			buildGeminiModelsURL(normalizedBaseURL),
		}
	}

	client := s.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	var (
		lastErr  error
		respBody []byte
	)

	for _, targetURL := range candidateURLs {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
		if reqErr != nil {
			lastErr = reqErr
			continue
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "sub2api-quick-sync/1.0")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("x-goog-api-key", apiKey)
			req.Header.Set("anthropic-version", "2023-06-01")
		}

		resp, doErr := client.Do(req)
		if doErr != nil {
			lastErr = doErr
			continue
		}

		bodyLimit := upstreamModelsBodyLimit
		if s.cfg != nil {
			bodyLimit = resolveModelsListReadLimit(s.cfg)
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, bodyLimit+1))
		_ = resp.Body.Close()

		if readErr != nil {
			lastErr = readErr
			continue
		}

		if resp.StatusCode == http.StatusOK {
			models, _, parseErr := extractUpstreamModelCatalog(body, false)
			if parseErr == nil && len(models) > 0 {
				respBody = body
				lastErr = nil
				break
			}
		}

		lastErr = fmt.Errorf("probe URL %s returned HTTP %d", targetURL, resp.StatusCode)
	}

	if len(respBody) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("failed to probe upstream models: %w", lastErr)
		}
		return nil, errors.New("upstream returned no models")
	}

	rawEntries, err := extractUpstreamModelRawEntries(respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upstream models: %w", err)
	}

	models, metadataMap, err := extractUpstreamModelCatalog(respBody, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upstream models catalog: %w", err)
	}

	if len(models) == 0 {
		return nil, errors.New("upstream returned no supported models")
	}

	// Build raw display name lookup map
	rawDisplayNameMap := make(map[string]string)
	for _, raw := range rawEntries {
		var entry quickSyncRawModelEntry
		if unmarshalErr := json.Unmarshal(raw, &entry); unmarshalErr == nil {
			entryID := NormalizeModelID(entry.ID)
			if entryID == "" {
				entryID = NormalizeModelID(entry.Slug)
			}
			if entryID == "" {
				entryID = NormalizeModelID(entry.Name)
			}

			displayName := strings.TrimSpace(entry.DisplayNameSnake)
			if displayName == "" {
				displayName = strings.TrimSpace(entry.DisplayNameCamel)
			}
			if displayName != "" {
				rawDisplayNameMap[entryID] = NormalizeModelDisplayName(displayName)
			}
		}
	}

	result := &QuickSyncProbeResult{
		Models: make([]QuickSyncModelItem, 0, len(models)),
		Total:  0,
	}

	for _, rawID := range models {
		normID := NormalizeModelID(rawID)
		if normID == "" {
			continue
		}

		displayName := normID
		if rawName, exists := rawDisplayNameMap[normID]; exists && rawName != "" {
			displayName = rawName
		} else if meta, ok := metadataMap[normID]; ok && strings.TrimSpace(meta.DisplayName) != "" {
			displayName = NormalizeModelDisplayName(meta.DisplayName)
		} else if meta, ok := metadataMap[rawID]; ok && strings.TrimSpace(meta.DisplayName) != "" {
			displayName = NormalizeModelDisplayName(meta.DisplayName)
		}
		if displayName == "" {
			displayName = normID
		}

		var basePriceIn, basePriceOut *float64

		// 1. Try PricingService (LiteLLM/models.dev)
		if s.pricingService != nil {
			if p := s.pricingService.GetModelPricing(normID); p != nil && !p.TokenPricingAbsent {
				basePriceIn = &p.InputCostPerToken
				basePriceOut = &p.OutputCostPerToken
			}
		}

		// 2. Try BillingService fallback
		if basePriceIn == nil && s.billingService != nil {
			if bp, err := s.billingService.GetModelPricing(normID); err == nil && bp != nil {
				basePriceIn = &bp.InputPricePerToken
				basePriceOut = &bp.OutputPricePerToken
			}
		}

		billingMode := string(BillingModeToken)
		var priceIn, priceOut *float64
		if basePriceIn != nil {
			val := *basePriceIn
			priceIn = &val
		}
		if basePriceOut != nil {
			val := *basePriceOut
			priceOut = &val
		}

		result.Models = append(result.Models, QuickSyncModelItem{
			ID:           normID,
			DisplayName:  displayName,
			BasePriceIn:  basePriceIn,
			BasePriceOut: basePriceOut,
			BillingMode:  billingMode,
			PriceIn:      priceIn,
			PriceOut:     priceOut,
		})
	}

	result.Total = len(result.Models)
	return result, nil
}

// ApplyBillingStrategy transforms probe model items by applying multiplier or per-request pricing rules.
func (s *ChannelQuickSyncService) ApplyBillingStrategy(items []QuickSyncModelItem, strategy QuickSyncBillingStrategy) []QuickSyncModelItem {
	out := make([]QuickSyncModelItem, len(items))
	for i, item := range items {
		cloned := item
		switch strategy.Mode {
		case string(BillingModePerRequest):
			cloned.BillingMode = string(BillingModePerRequest)
			if strategy.PerRequestPrice != nil {
				val := *strategy.PerRequestPrice
				cloned.PerReqPrice = &val
			}
			cloned.PriceIn = nil
			cloned.PriceOut = nil
		case "ratio":
			cloned.BillingMode = string(BillingModeToken)
			ratio := strategy.Ratio
			if ratio <= 0 {
				ratio = 1.0
			}
			if item.BasePriceIn != nil {
				val := *item.BasePriceIn * ratio
				cloned.PriceIn = &val
			}
			if item.BasePriceOut != nil {
				val := *item.BasePriceOut * ratio
				cloned.PriceOut = &val
			}
			cloned.PerReqPrice = nil
		case "fixed":
			// Leave current price overrides intact
		default:
			cloned.BillingMode = string(BillingModeToken)
		}
		out[i] = cloned
	}
	return out
}

// CommitQuickSync atomically creates an account, optionally creates a new group, binds groups and model allowlists, creates a channel, and configures channel model pricing.
func (s *ChannelQuickSyncService) CommitQuickSync(ctx context.Context, params QuickSyncCommitParams) (*QuickSyncCommitResult, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	baseURL := strings.TrimSpace(params.BaseURL)
	if baseURL == "" {
		return nil, errors.New("base_url is required")
	}

	apiKey := strings.TrimSpace(params.APIKey)
	if apiKey == "" {
		return nil, errors.New("api_key is required")
	}

	if len(params.Models) == 0 {
		return nil, errors.New("models cannot be empty")
	}

	platform := strings.ToLower(strings.TrimSpace(params.Platform))
	if platform == "" {
		platform = PlatformAntigravity
	}

	normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base_url: %w", err)
	}

	if s.accountRepo == nil || s.groupRepo == nil || s.channelRepo == nil {
		return nil, errors.New("repositories not configured")
	}

	var (
		tx    *dbent.Tx
		txErr error
	)
	opCtx := ctx
	if s.entClient != nil {
		tx, txErr = s.entClient.Tx(ctx)
		if txErr != nil {
			return nil, fmt.Errorf("begin transaction: %w", txErr)
		}
		defer func() { _ = tx.Rollback() }()
		opCtx = dbent.NewTxContext(ctx, tx)
	}

	// 1. Optional new group creation
	var newGroupID int64
	if params.NewGroup != nil && params.NewGroup.Create {
		groupName := strings.TrimSpace(params.NewGroup.Name)
		if groupName == "" {
			groupName = name
		}
		rateMultiplier := params.NewGroup.RateMultiplier
		if rateMultiplier <= 0 {
			rateMultiplier = 1.0
		}
		exists, err := s.groupRepo.ExistsByName(opCtx, groupName)
		if err != nil {
			return nil, fmt.Errorf("check group exists: %w", err)
		}
		if exists {
			return nil, ErrGroupExists
		}
		newGroup := &Group{
			Name:           groupName,
			Description:    fmt.Sprintf("QuickSync group for %s", name),
			Platform:       platform,
			RateMultiplier: rateMultiplier,
			Status:         StatusActive,
		}
		if err := s.groupRepo.Create(opCtx, newGroup); err != nil {
			return nil, fmt.Errorf("create new group: %w", err)
		}
		newGroupID = newGroup.ID
	}

	// 2. Resolve model assignments and validate pricing
	type resolvedModelItem struct {
		model           string
		targetGroupID   int64
		billingMode     string
		inputPrice      *float64
		outputPrice     *float64
		perRequestPrice *float64
	}

	resolvedModels := make([]resolvedModelItem, 0, len(params.Models))
	groupToModelsMap := make(map[int64][]string)
	distinctGroupIDsMap := make(map[int64]bool)
	seenModels := make(map[string]bool)

	for _, item := range params.Models {
		normModel := NormalizeModelID(item.Model)
		if normModel == "" {
			return nil, errors.New("model name cannot be empty")
		}
		lowerModel := strings.ToLower(normModel)
		if seenModels[lowerModel] {
			return nil, fmt.Errorf("duplicate model %s in models list", normModel)
		}
		seenModels[lowerModel] = true

		targetGID := item.TargetGroupID
		if targetGID <= 0 {
			if newGroupID > 0 {
				targetGID = newGroupID
			} else if params.DefaultGroupID != nil && *params.DefaultGroupID > 0 {
				targetGID = *params.DefaultGroupID
			}
		}
		if targetGID <= 0 {
			return nil, fmt.Errorf("target group not specified for model %s", normModel)
		}

		bMode := strings.ToLower(strings.TrimSpace(item.BillingMode))
		if bMode == "" {
			if params.BillingStrategy != nil && params.BillingStrategy.Mode == string(BillingModePerRequest) {
				bMode = string(BillingModePerRequest)
			} else {
				bMode = string(BillingModeToken)
			}
		}

		inputPrice := item.InputPrice
		outputPrice := item.OutputPrice
		perReqPrice := item.PerRequestPrice

		if bMode == string(BillingModePerRequest) {
			if perReqPrice == nil && params.BillingStrategy != nil && params.BillingStrategy.PerRequestPrice != nil {
				perReqPrice = params.BillingStrategy.PerRequestPrice
			}
			if perReqPrice == nil {
				return nil, fmt.Errorf("per_request_price required for model %s with per_request billing mode", normModel)
			}
			if *perReqPrice < 0 {
				return nil, fmt.Errorf("per_request_price must be >= 0 for model %s", normModel)
			}
		} else {
			bMode = string(BillingModeToken)
			if inputPrice != nil && *inputPrice < 0 {
				return nil, fmt.Errorf("input_price must be >= 0 for model %s", normModel)
			}
			if outputPrice != nil && *outputPrice < 0 {
				return nil, fmt.Errorf("output_price must be >= 0 for model %s", normModel)
			}
		}

		resolvedModels = append(resolvedModels, resolvedModelItem{
			model:           normModel,
			targetGroupID:   targetGID,
			billingMode:     bMode,
			inputPrice:      inputPrice,
			outputPrice:     outputPrice,
			perRequestPrice: perReqPrice,
		})
		groupToModelsMap[targetGID] = append(groupToModelsMap[targetGID], normModel)
		distinctGroupIDsMap[targetGID] = true
	}

	distinctGroupIDs := make([]int64, 0, len(distinctGroupIDsMap))
	for gid := range distinctGroupIDsMap {
		distinctGroupIDs = append(distinctGroupIDs, gid)
	}
	sort.Slice(distinctGroupIDs, func(i, j int) bool {
		return distinctGroupIDs[i] < distinctGroupIDs[j]
	})

	// 3. Validate groups & append ModelAllowlist
	for _, gid := range distinctGroupIDs {
		if gid == newGroupID {
			continue
		}
		g, err := s.groupRepo.GetByID(opCtx, gid)
		if err != nil {
			return nil, fmt.Errorf("get group %d: %w", gid, err)
		}
		if g.RequireOAuthOnly && (g.Platform == PlatformOpenAI || g.Platform == PlatformAntigravity || g.Platform == PlatformAnthropic || g.Platform == PlatformGemini || g.Platform == PlatformGrok) {
			return nil, fmt.Errorf("group [%s] requires OAuth accounts, cannot bind api_key account", g.Name)
		}

		existingMap := make(map[string]bool)
		for _, m := range g.ModelAllowlist.Models {
			existingMap[strings.ToLower(m)] = true
		}
		changed := false
		for _, m := range groupToModelsMap[gid] {
			if !existingMap[strings.ToLower(m)] {
				g.ModelAllowlist.Models = append(g.ModelAllowlist.Models, m)
				existingMap[strings.ToLower(m)] = true
				changed = true
			}
		}
		if changed {
			if err := s.groupRepo.Update(opCtx, g); err != nil {
				return nil, fmt.Errorf("update group %d model allowlist: %w", gid, err)
			}
		}
	}

	if newGroupID > 0 {
		newGroup, err := s.groupRepo.GetByID(opCtx, newGroupID)
		if err == nil && newGroup != nil {
			newGroup.ModelAllowlist.Models = groupToModelsMap[newGroupID]
			_ = s.groupRepo.Update(opCtx, newGroup)
		}
	}

	// 4. Create Account and bind groups
	account := &Account{
		Name:     name,
		Platform: platform,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  apiKey,
			"base_url": normalizedBaseURL,
		},
		Status:             StatusActive,
		AutoPauseOnExpired: true,
	}
	if err := s.accountRepo.Create(opCtx, account); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	if len(distinctGroupIDs) > 0 {
		if err := s.accountRepo.BindGroups(opCtx, account.ID, distinctGroupIDs); err != nil {
			return nil, fmt.Errorf("bind account groups: %w", err)
		}
	}

	// 5. Check channel conflicts and create channel with model pricing
	channelExists, err := s.channelRepo.ExistsByName(opCtx, name)
	if err != nil {
		return nil, fmt.Errorf("check channel exists: %w", err)
	}
	if channelExists {
		return nil, ErrChannelExists
	}

	conflictingGroups, err := s.channelRepo.GetGroupsInOtherChannels(opCtx, 0, distinctGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("check group conflicts: %w", err)
	}
	if len(conflictingGroups) > 0 {
		return nil, ErrGroupAlreadyInChannel
	}

	pricingList := make([]ChannelModelPricing, 0, len(resolvedModels))
	for _, item := range resolvedModels {
		var p ChannelModelPricing
		p.Platform = platform
		p.Models = []string{item.model}
		if item.billingMode == string(BillingModePerRequest) {
			zero := 0.0
			p.BillingMode = BillingModePerRequest
			p.InputPrice = &zero
			p.OutputPrice = &zero
			p.PerRequestPrice = item.perRequestPrice
		} else {
			p.BillingMode = BillingModeToken
			p.InputPrice = item.inputPrice
			p.OutputPrice = item.outputPrice
		}
		pricingList = append(pricingList, p)
	}

	channel := &Channel{
		Name:               name,
		Description:        fmt.Sprintf("QuickSync channel for %s", name),
		Status:             StatusActive,
		BillingModelSource: BillingModelSourceChannelMapped,
		GroupIDs:           distinctGroupIDs,
		ModelPricing:       pricingList,
	}

	if err := s.channelRepo.Create(opCtx, channel); err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}

	// 6. Commit transaction
	if tx != nil {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit transaction: %w", err)
		}
	}

	// 7. Invalidate channel cache
	if s.channelCacheInvalidator != nil {
		s.channelCacheInvalidator.InvalidateCache()
	} else if s.channelService != nil {
		s.channelService.InvalidateCache()
	} else if s.cachePubSub != nil {
		_ = s.cachePubSub.NotifyUpdate(ctx)
	}

	return &QuickSyncCommitResult{
		ChannelID:  channel.ID,
		AccountID:  account.ID,
		GroupIDs:   distinctGroupIDs,
		ModelCount: len(resolvedModels),
	}, nil
}
