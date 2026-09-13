package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// ForwardUpstream 使用 base_url + /v1/messages + 双 header 认证透传上游 Claude 请求
func (s *AntigravityGatewayService) ForwardUpstream(ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
	beginUpstreamResponseModelObservation(c)
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	// 获取上游配置
	baseURL := strings.TrimSpace(account.GetCredential("base_url"))
	apiKey := strings.TrimSpace(account.GetCredential("api_key"))
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("upstream account missing base_url or api_key")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	// 解析请求获取模型信息
	var claudeReq antigravity.ClaudeRequest
	if err := json.Unmarshal(body, &claudeReq); err != nil {
		return nil, fmt.Errorf("parse claude request: %w", err)
	}
	if strings.TrimSpace(claudeReq.Model) == "" {
		return nil, fmt.Errorf("missing model")
	}
	originalModel := claudeReq.Model

	// 构建上游请求 URL
	upstreamURL := baseURL + "/v1/messages"

	// 能力维度 sanitize：Anthropic-compatible 上游透传路径也需要保证 body↔beta header
	// 对称。客户端 anthropic-beta header 不含 context-management-2025-06-27 但 body 带
	// context_management 时 strip，与 Anthropic 直连 / Bedrock / Vertex 路径保持一致。
	clientBeta := c.GetHeader("anthropic-beta")
	if sanitized, changed := sanitizeAnthropicBodyForBetaTokens(body, clientBeta); changed {
		body = sanitized
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("x-api-key", apiKey) // Claude API 兼容

	// 透传 Claude 相关 headers
	if v := c.GetHeader("anthropic-version"); v != "" {
		req.Header.Set("anthropic-version", v)
	}
	if v := clientBeta; v != "" {
		req.Header.Set("anthropic-beta", v)
	}

	// 代理 URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	// 发送请求
	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s upstream request failed: %v", prefix, err)
		return nil, fmt.Errorf("upstream request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 处理错误响应
	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)

		// 429 错误时标记账号限流
		if resp.StatusCode == http.StatusTooManyRequests {
			s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody, originalModel, 0, "", false)
		}

		// 透传上游错误
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(resp.StatusCode)
		_, _ = c.Writer.Write(respBody)

		return &ForwardResult{
			Model: originalModel,
		}, nil
	}

	// 处理成功响应（流式/非流式）
	var usage *ClaudeUsage
	var firstTokenMs *int
	var clientDisconnect bool

	if claudeReq.Stream {
		// 流式响应：透传
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		c.Status(http.StatusOK)

		streamRes := s.streamUpstreamResponse(c, resp, startTime)
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
		clientDisconnect = streamRes.clientDisconnect
	} else {
		// 非流式响应：直接透传
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read upstream response: %w", err)
		}

		// 提取 usage
		upstreamResponseModelObserverFromContext(c).ObserveAnthropic(respBody)
		usage = s.extractClaudeUsage(respBody)

		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write(respBody)
	}

	// 构建计费结果
	duration := time.Since(startTime)
	logger.LegacyPrintf("service.antigravity_gateway", "%s status=success duration_ms=%d", prefix, duration.Milliseconds())

	return &ForwardResult{
		Model:                         originalModel,
		UpstreamResponseModel:         observedUpstreamResponseModel(c),
		UpstreamResponseModelConflict: observedUpstreamResponseModelConflict(c),
		Stream:                        claudeReq.Stream,
		Duration:                      duration,
		FirstTokenMs:                  firstTokenMs,
		ClientDisconnect:              clientDisconnect,
		UpstreamHeaders:               resp.Header,
		Usage: ClaudeUsage{
			InputTokens:              usage.InputTokens,
			OutputTokens:             usage.OutputTokens,
			CacheReadInputTokens:     usage.CacheReadInputTokens,
			CacheCreationInputTokens: usage.CacheCreationInputTokens,
		},
	}, nil
}

// streamUpstreamResponse 透传上游 SSE 流并提取 Claude usage
func (s *AntigravityGatewayService) streamUpstreamResponse(c *gin.Context, resp *http.Response, startTime time.Time) *antigravityStreamResult {
	return s.streamUpstreamGeneric(c, resp, startTime, "antigravity upstream", "event: ping\ndata: {\"type\": \"ping\"}\n\n", func(line string, usage *ClaudeUsage) {
		if data, ok := extractAnthropicSSEDataLine(line); ok {
			upstreamResponseModelObserverFromContext(c).ObserveAnthropic([]byte(strings.TrimSpace(data)))
		}
		s.extractSSEUsage(line, usage)
	})
}

// streamUpstreamGeneric 透传上游 SSE 流通用实现
func (s *AntigravityGatewayService) streamUpstreamGeneric(
	c *gin.Context,
	resp *http.Response,
	startTime time.Time,
	logTag string,
	pingMsg string,
	handler func(line string, usage *ClaudeUsage),
) *antigravityStreamResult {
	usage := &ClaudeUsage{}
	var firstTokenMs *int

	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.settingService.cfg.Gateway.MaxLineSize
	}
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	type scanEvent struct {
		line string
		err  error
	}
	events := make(chan scanEvent, 16)
	done := make(chan struct{})
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
		}
	}
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()}) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: err})
		}
	}()
	defer close(done)

	streamInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.settingService.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
	}
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
	}
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
	}

	keepaliveInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamKeepaliveInterval > 0 {
		keepaliveInterval = time.Duration(s.settingService.cfg.Gateway.StreamKeepaliveInterval) * time.Second
	}
	var keepaliveTicker *time.Ticker
	if keepaliveInterval > 0 {
		keepaliveTicker = time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()
	}
	var keepaliveCh <-chan time.Time
	if keepaliveTicker != nil {
		keepaliveCh = keepaliveTicker.C
	}
	lastDataAt := time.Now()

	flusher, _ := c.Writer.(http.Flusher)
	cw := newAntigravityClientWriter(c.Writer, flusher, logTag)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMs, clientDisconnect: cw.Disconnected()}
			}
			if ev.err != nil {
				if disconnect, handled := handleStreamReadError(ev.err, cw.Disconnected(), logTag); handled {
					return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMs, clientDisconnect: disconnect}
				}
				logger.LegacyPrintf("service.antigravity_gateway", "Stream read error (%s): %v", logTag, ev.err)
				return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMs}
			}

			lastDataAt = time.Now()
			line := ev.line

			if firstTokenMs == nil && len(line) > 0 {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
			}

			if handler != nil {
				handler(line, usage)
			}

			cw.Fprintf("%s\n", line)

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
			}
			if cw.Disconnected() {
				logger.LegacyPrintf("service.antigravity_gateway", "Upstream timeout after client disconnect (%s), returning collected usage", logTag)
				return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMs, clientDisconnect: true}
			}
			logger.LegacyPrintf("service.antigravity_gateway", "Stream data interval timeout (%s)", logTag)
			return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMs}

		case <-keepaliveCh:
			if cw.Disconnected() {
				continue
			}
			if time.Since(lastDataAt) < keepaliveInterval {
				continue
			}
			msg := pingMsg
			if msg == "" {
				msg = ": ping\n\n"
			}
			if !cw.Fprintf("%s", msg) {
				logger.LegacyPrintf("service.antigravity_gateway", "Client disconnected during keepalive ping (%s), continuing to drain upstream for billing", logTag)
				continue
			}
		}
	}
}

// extractSSEUsage 从 SSE data 行中提取 Claude usage（用于流式透传场景）
//
// Anthropic streaming 的 usage 字段分布在两类事件中：
//   - message_start：嵌套在 event.message.usage（input_tokens、cache_creation_input_tokens、
//     cache_read_input_tokens 等输入侧字段）
//   - message_delta：位于顶层 event.usage（流结束时的最终 output_tokens）
//
// 仅读取顶层 event.usage 会漏掉 message_start 的输入侧字段，导致流式透传请求落库的
// usage_logs 记录 input_tokens=0。
func (s *AntigravityGatewayService) extractSSEUsage(line string, usage *ClaudeUsage) {
	if !strings.HasPrefix(line, "data: ") {
		return
	}
	dataStr := strings.TrimPrefix(line, "data: ")
	var event map[string]any
	if json.Unmarshal([]byte(dataStr), &event) != nil {
		return
	}
	var u map[string]any
	if eventType, _ := event["type"].(string); eventType == "message_start" {
		if msg, ok := event["message"].(map[string]any); ok {
			u, _ = msg["usage"].(map[string]any)
		}
	} else {
		u, _ = event["usage"].(map[string]any)
	}
	if u == nil {
		return
	}
	if v, ok := u["input_tokens"].(float64); ok && int(v) > 0 {
		usage.InputTokens = int(v)
	}
	if v, ok := u["output_tokens"].(float64); ok && int(v) > 0 {
		usage.OutputTokens = int(v)
	}
	if v, ok := u["cache_read_input_tokens"].(float64); ok && int(v) > 0 {
		usage.CacheReadInputTokens = int(v)
	}
	if v, ok := u["cache_creation_input_tokens"].(float64); ok && int(v) > 0 {
		usage.CacheCreationInputTokens = int(v)
	}
	// 解析嵌套的 cache_creation 对象中的 5m/1h 明细
	if cc, ok := u["cache_creation"].(map[string]any); ok {
		if v, ok := cc["ephemeral_5m_input_tokens"].(float64); ok {
			usage.CacheCreation5mTokens = int(v)
		}
		if v, ok := cc["ephemeral_1h_input_tokens"].(float64); ok {
			usage.CacheCreation1hTokens = int(v)
		}
	}
}

// extractClaudeUsage 从非流式 Claude 响应提取 usage
func (s *AntigravityGatewayService) extractClaudeUsage(body []byte) *ClaudeUsage {
	usage := &ClaudeUsage{}
	var resp map[string]any
	if json.Unmarshal(body, &resp) != nil {
		return usage
	}
	if u, ok := resp["usage"].(map[string]any); ok {
		if v, ok := u["input_tokens"].(float64); ok {
			usage.InputTokens = int(v)
		}
		if v, ok := u["output_tokens"].(float64); ok {
			usage.OutputTokens = int(v)
		}
		if v, ok := u["cache_read_input_tokens"].(float64); ok {
			usage.CacheReadInputTokens = int(v)
		}
		if v, ok := u["cache_creation_input_tokens"].(float64); ok {
			usage.CacheCreationInputTokens = int(v)
		}
		// 解析嵌套的 cache_creation 对象中的 5m/1h 明细
		if cc, ok := u["cache_creation"].(map[string]any); ok {
			if v, ok := cc["ephemeral_5m_input_tokens"].(float64); ok {
				usage.CacheCreation5mTokens = int(v)
			}
			if v, ok := cc["ephemeral_1h_input_tokens"].(float64); ok {
				usage.CacheCreation1hTokens = int(v)
			}
		}
	}
	return usage
}

// ForwardUpstreamChatCompletions forwards OpenAI Chat Completions requests to the upstream base_url + /v1/chat/completions
func (s *AntigravityGatewayService) ForwardUpstreamChatCompletions(ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
	beginUpstreamResponseModelObservation(c)
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	baseURL := strings.TrimSpace(account.GetCredential("base_url"))
	apiKey := strings.TrimSpace(account.GetCredential("api_key"))
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("upstream account missing base_url or api_key")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	model := gjson.GetBytes(body, "model").String()
	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("missing model")
	}
	stream := gjson.GetBytes(body, "stream").Bool()

	upstreamURL := baseURL + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s chat completions upstream request failed: %v", prefix, err)
		return nil, fmt.Errorf("upstream request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)
		if resp.StatusCode == http.StatusTooManyRequests {
			s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody, model, 0, "", false)
		}
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(resp.StatusCode)
		_, _ = c.Writer.Write(respBody)
		return &ForwardResult{
			Model: model,
		}, nil
	}

	var usage *ClaudeUsage
	var firstTokenMs *int
	var clientDisconnect bool

	if stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		c.Status(http.StatusOK)

		streamRes := s.streamUpstreamGeneric(c, resp, startTime, "antigravity cc upstream", ": ping\n\n", func(line string, u *ClaudeUsage) {
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if strings.TrimSpace(data) != "[DONE]" {
					upstreamResponseModelObserverFromContext(c).ObserveOpenAI([]byte(data), "")
					s.extractOpenAISSEUsage(line, u)
				}
			}
		})
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
		clientDisconnect = streamRes.clientDisconnect
	} else {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read upstream response: %w", err)
		}
		upstreamResponseModelObserverFromContext(c).ObserveOpenAI(respBody, "")
		usage = s.extractOpenAIUsage(respBody)

		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write(respBody)
	}

	duration := time.Since(startTime)
	logger.LegacyPrintf("service.antigravity_gateway", "%s chat completions status=success duration_ms=%d", prefix, duration.Milliseconds())

	finalUsage := ClaudeUsage{}
	if usage != nil {
		finalUsage = *usage
	}

	return &ForwardResult{
		Model:                         model,
		UpstreamResponseModel:         observedUpstreamResponseModel(c),
		UpstreamResponseModelConflict: observedUpstreamResponseModelConflict(c),
		Stream:                        stream,
		Duration:                      duration,
		FirstTokenMs:                  firstTokenMs,
		ClientDisconnect:              clientDisconnect,
		UpstreamHeaders:               resp.Header,
		Usage:                         finalUsage,
	}, nil
}

// ForwardUpstreamResponses forwards OpenAI Responses requests to the upstream base_url + /v1/responses
func (s *AntigravityGatewayService) ForwardUpstreamResponses(ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
	beginUpstreamResponseModelObservation(c)
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	baseURL := strings.TrimSpace(account.GetCredential("base_url"))
	apiKey := strings.TrimSpace(account.GetCredential("api_key"))
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("upstream account missing base_url or api_key")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	model := gjson.GetBytes(body, "model").String()
	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("missing model")
	}
	stream := gjson.GetBytes(body, "stream").Bool()

	upstreamURL := baseURL + "/v1/responses"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s responses upstream request failed: %v", prefix, err)
		return nil, fmt.Errorf("upstream request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)
		if resp.StatusCode == http.StatusTooManyRequests {
			s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody, model, 0, "", false)
		}
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(resp.StatusCode)
		_, _ = c.Writer.Write(respBody)
		return &ForwardResult{
			Model: model,
		}, nil
	}

	var usage *ClaudeUsage
	var firstTokenMs *int
	var clientDisconnect bool

	if stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		c.Status(http.StatusOK)

		streamRes := s.streamUpstreamGeneric(c, resp, startTime, "antigravity responses upstream", ": ping\n\n", func(line string, u *ClaudeUsage) {
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if strings.TrimSpace(data) != "[DONE]" {
					upstreamResponseModelObserverFromContext(c).ObserveOpenAI([]byte(data), "")
					s.extractResponsesSSEUsage(line, u)
				}
			}
		})
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
		clientDisconnect = streamRes.clientDisconnect
	} else {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read upstream response: %w", err)
		}
		upstreamResponseModelObserverFromContext(c).ObserveOpenAI(respBody, "")
		usage = s.extractResponsesUsage(respBody)

		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write(respBody)
	}

	duration := time.Since(startTime)
	logger.LegacyPrintf("service.antigravity_gateway", "%s responses status=success duration_ms=%d", prefix, duration.Milliseconds())

	finalUsage := ClaudeUsage{}
	if usage != nil {
		finalUsage = *usage
	}

	return &ForwardResult{
		Model:                         model,
		UpstreamResponseModel:         observedUpstreamResponseModel(c),
		UpstreamResponseModelConflict: observedUpstreamResponseModelConflict(c),
		Stream:                        stream,
		Duration:                      duration,
		FirstTokenMs:                  firstTokenMs,
		ClientDisconnect:              clientDisconnect,
		UpstreamHeaders:               resp.Header,
		Usage:                         finalUsage,
	}, nil
}

// ForwardUpstreamGemini forwards Gemini native requests to the upstream base_url + /v1beta/models/...
func (s *AntigravityGatewayService) ForwardUpstreamGemini(ctx context.Context, c *gin.Context, account *Account, modelName string, action string, stream bool, body []byte) (*ForwardResult, error) {
	beginUpstreamResponseModelObservation(c)
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	baseURL := strings.TrimSpace(account.GetCredential("base_url"))
	apiKey := strings.TrimSpace(account.GetCredential("api_key"))
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("upstream account missing base_url or api_key")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	upstreamURL := fmt.Sprintf("%s/v1beta/models/%s:%s", baseURL, modelName, action)
	if stream {
		upstreamURL += "?alt=sse"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("x-goog-api-key", apiKey)

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s gemini upstream request failed: %v", prefix, err)
		return nil, fmt.Errorf("upstream request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)
		if resp.StatusCode == http.StatusTooManyRequests {
			s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody, modelName, 0, "", false)
		}
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(resp.StatusCode)
		_, _ = c.Writer.Write(respBody)
		return &ForwardResult{
			Model: modelName,
		}, nil
	}

	var usage *ClaudeUsage
	var firstTokenMs *int
	var clientDisconnect bool

	if stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		c.Status(http.StatusOK)

		streamRes := s.streamUpstreamGeneric(c, resp, startTime, "antigravity gemini upstream", ": ping\n\n", func(line string, u *ClaudeUsage) {
			s.extractGeminiSSEUsage(line, u)
		})
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
		clientDisconnect = streamRes.clientDisconnect
	} else {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read upstream response: %w", err)
		}
		usage = s.extractGeminiUsage(respBody)

		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write(respBody)
	}

	duration := time.Since(startTime)
	logger.LegacyPrintf("service.antigravity_gateway", "%s gemini status=success duration_ms=%d", prefix, duration.Milliseconds())

	finalUsage := ClaudeUsage{}
	if usage != nil {
		finalUsage = *usage
	}

	return &ForwardResult{
		Model:                         modelName,
		UpstreamResponseModel:         observedUpstreamResponseModel(c),
		UpstreamResponseModelConflict: observedUpstreamResponseModelConflict(c),
		Stream:                        stream,
		Duration:                      duration,
		FirstTokenMs:                  firstTokenMs,
		ClientDisconnect:              clientDisconnect,
		UpstreamHeaders:               resp.Header,
		Usage:                         finalUsage,
	}, nil
}

func (s *AntigravityGatewayService) extractOpenAIUsage(body []byte) *ClaudeUsage {
	usage := &ClaudeUsage{}
	res := gjson.GetBytes(body, "usage")
	if !res.Exists() {
		return usage
	}
	usage.InputTokens = int(res.Get("prompt_tokens").Int())
	usage.OutputTokens = int(res.Get("completion_tokens").Int())
	usage.CacheReadInputTokens = int(res.Get("prompt_tokens_details.cached_tokens").Int())
	return usage
}

func (s *AntigravityGatewayService) extractOpenAISSEUsage(line string, usage *ClaudeUsage) {
	if !strings.HasPrefix(line, "data: ") {
		return
	}
	data := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if data == "[DONE]" || data == "" {
		return
	}
	res := gjson.Get(data, "usage")
	if !res.Exists() {
		return
	}
	if in := res.Get("prompt_tokens").Int(); in > 0 {
		usage.InputTokens = int(in)
	}
	if out := res.Get("completion_tokens").Int(); out > 0 {
		usage.OutputTokens = int(out)
	}
	if cached := res.Get("prompt_tokens_details.cached_tokens").Int(); cached > 0 {
		usage.CacheReadInputTokens = int(cached)
	}
}

func (s *AntigravityGatewayService) extractResponsesUsage(body []byte) *ClaudeUsage {
	usage := &ClaudeUsage{}
	res := gjson.GetBytes(body, "usage")
	if !res.Exists() {
		res = gjson.GetBytes(body, "response.usage")
	}
	if !res.Exists() {
		return usage
	}
	usage.InputTokens = int(res.Get("input_tokens").Int())
	usage.OutputTokens = int(res.Get("output_tokens").Int())
	usage.CacheReadInputTokens = int(res.Get("input_token_details.cached_tokens").Int())
	return usage
}

func (s *AntigravityGatewayService) extractResponsesSSEUsage(line string, usage *ClaudeUsage) {
	if !strings.HasPrefix(line, "data: ") {
		return
	}
	data := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if data == "[DONE]" || data == "" {
		return
	}
	res := gjson.Get(data, "usage")
	if !res.Exists() {
		res = gjson.Get(data, "response.usage")
	}
	if !res.Exists() {
		return
	}
	if in := res.Get("input_tokens").Int(); in > 0 {
		usage.InputTokens = int(in)
	}
	if out := res.Get("output_tokens").Int(); out > 0 {
		usage.OutputTokens = int(out)
	}
	if cached := res.Get("input_token_details.cached_tokens").Int(); cached > 0 {
		usage.CacheReadInputTokens = int(cached)
	}
}

func (s *AntigravityGatewayService) extractGeminiUsage(body []byte) *ClaudeUsage {
	usage := &ClaudeUsage{}
	res := gjson.GetBytes(body, "usageMetadata")
	if !res.Exists() {
		return usage
	}
	usage.InputTokens = int(res.Get("promptTokenCount").Int())
	usage.OutputTokens = int(res.Get("candidatesTokenCount").Int())
	usage.CacheReadInputTokens = int(res.Get("cachedContentTokenCount").Int())
	return usage
}

func (s *AntigravityGatewayService) extractGeminiSSEUsage(line string, usage *ClaudeUsage) {
	if !strings.HasPrefix(line, "data: ") {
		return
	}
	data := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if data == "" {
		return
	}
	res := gjson.Get(data, "usageMetadata")
	if !res.Exists() {
		return
	}
	if in := res.Get("promptTokenCount").Int(); in > 0 {
		usage.InputTokens = int(in)
	}
	if out := res.Get("candidatesTokenCount").Int(); out > 0 {
		usage.OutputTokens = int(out)
	}
	if cached := res.Get("cachedContentTokenCount").Int(); cached > 0 {
		usage.CacheReadInputTokens = int(cached)
	}
}
