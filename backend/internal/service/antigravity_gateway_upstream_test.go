package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestExtractOpenAIUsage(t *testing.T) {
	svc := &AntigravityGatewayService{}
	body := []byte(`{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"usage": {
			"prompt_tokens": 15,
			"completion_tokens": 25,
			"total_tokens": 40,
			"prompt_tokens_details": {
				"cached_tokens": 5
			}
		}
	}`)
	usage := svc.extractOpenAIUsage(body)
	require.Equal(t, 15, usage.InputTokens)
	require.Equal(t, 25, usage.OutputTokens)
	require.Equal(t, 5, usage.CacheReadInputTokens)
}

func TestExtractOpenAISSEUsage(t *testing.T) {
	svc := &AntigravityGatewayService{}
	usage := &ClaudeUsage{}
	svc.extractOpenAISSEUsage("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":20,\"prompt_tokens_details\":{\"cached_tokens\":3}}}", usage)
	require.Equal(t, 10, usage.InputTokens)
	require.Equal(t, 20, usage.OutputTokens)
	require.Equal(t, 3, usage.CacheReadInputTokens)
}

func TestExtractResponsesUsage(t *testing.T) {
	svc := &AntigravityGatewayService{}
	body := []byte(`{
		"id": "resp-123",
		"usage": {
			"input_tokens": 30,
			"output_tokens": 50,
			"input_token_details": {
				"cached_tokens": 12
			}
		}
	}`)
	usage := svc.extractResponsesUsage(body)
	require.Equal(t, 30, usage.InputTokens)
	require.Equal(t, 50, usage.OutputTokens)
	require.Equal(t, 12, usage.CacheReadInputTokens)
}

func TestExtractResponsesSSEUsage(t *testing.T) {
	svc := &AntigravityGatewayService{}
	usage := &ClaudeUsage{}
	svc.extractResponsesSSEUsage("data: {\"response\":{\"usage\":{\"input_tokens\":8,\"output_tokens\":16,\"input_token_details\":{\"cached_tokens\":2}}}}", usage)
	require.Equal(t, 8, usage.InputTokens)
	require.Equal(t, 16, usage.OutputTokens)
	require.Equal(t, 2, usage.CacheReadInputTokens)
}

func TestExtractGeminiUsageAntigravityUpstream(t *testing.T) {
	svc := &AntigravityGatewayService{}
	body := []byte(`{
		"candidates": [{"content": {"parts": [{"text": "hello"}]}}],
		"usageMetadata": {
			"promptTokenCount": 100,
			"candidatesTokenCount": 200,
			"cachedContentTokenCount": 40
		}
	}`)
	usage := svc.extractGeminiUsage(body)
	require.Equal(t, 100, usage.InputTokens)
	require.Equal(t, 200, usage.OutputTokens)
	require.Equal(t, 40, usage.CacheReadInputTokens)
}

func TestExtractGeminiSSEUsage(t *testing.T) {
	svc := &AntigravityGatewayService{}
	usage := &ClaudeUsage{}
	svc.extractGeminiSSEUsage("data: {\"candidates\":[],\"usageMetadata\":{\"promptTokenCount\":45,\"candidatesTokenCount\":90,\"cachedContentTokenCount\":15}}", usage)
	require.Equal(t, 45, usage.InputTokens)
	require.Equal(t, 90, usage.OutputTokens)
	require.Equal(t, 15, usage.CacheReadInputTokens)
}

func TestForwardUpstreamChatCompletions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	var capturedReq *http.Request
	upstreamResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(`{"id":"chatcmpl-1","choices":[{"message":{"role":"assistant","content":"hello"}}],"usage":{"prompt_tokens":10,"completion_tokens":20}}`)),
	}

	stub := &httpUpstreamStubFunc{
		doFunc: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
			capturedReq = req
			return upstreamResp, nil
		},
	}

	svc := &AntigravityGatewayService{
		httpUpstream: stub,
	}

	account := &Account{
		ID:       101,
		Name:     "upstream-node",
		Platform: PlatformAntigravity,
		Type:     AccountTypeUpstream,
		Credentials: map[string]any{
			"base_url": "http://154.36.173.146",
			"api_key":  "test-secret-key",
		},
	}

	body := []byte(`{"model":"gemini-3-flash","messages":[{"role":"user","content":"hi"}]}`)
	res, err := svc.ForwardUpstreamChatCompletions(context.Background(), c, account, body)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "gemini-3-flash", res.Model)
	require.Equal(t, 10, res.Usage.InputTokens)
	require.Equal(t, 20, res.Usage.OutputTokens)
	require.Equal(t, "http://154.36.173.146/v1/chat/completions", capturedReq.URL.String())
	require.Equal(t, "Bearer test-secret-key", capturedReq.Header.Get("Authorization"))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestForwardUpstreamResponses_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

	var capturedReq *http.Request
	upstreamResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(`{"id":"resp-1","output":[],"usage":{"input_tokens":12,"output_tokens":24}}`)),
	}

	stub := &httpUpstreamStubFunc{
		doFunc: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
			capturedReq = req
			return upstreamResp, nil
		},
	}

	svc := &AntigravityGatewayService{
		httpUpstream: stub,
	}

	account := &Account{
		ID:       102,
		Name:     "upstream-node-resp",
		Platform: PlatformAntigravity,
		Type:     AccountTypeUpstream,
		Credentials: map[string]any{
			"base_url": "http://154.36.173.146",
			"api_key":  "test-secret-key",
		},
	}

	body := []byte(`{"model":"gemini-3-flash","input":"hi"}`)
	res, err := svc.ForwardUpstreamResponses(context.Background(), c, account, body)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "gemini-3-flash", res.Model)
	require.Equal(t, 12, res.Usage.InputTokens)
	require.Equal(t, 24, res.Usage.OutputTokens)
	require.Equal(t, "http://154.36.173.146/v1/responses", capturedReq.URL.String())
	require.Equal(t, "Bearer test-secret-key", capturedReq.Header.Get("Authorization"))
}

func TestForwardUpstreamGemini_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-3-flash:generateContent", nil)

	var capturedReq *http.Request
	upstreamResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(`{"candidates":[{"content":{"parts":[{"text":"hello"}]}}],"usageMetadata":{"promptTokenCount":50,"candidatesTokenCount":100}}`)),
	}

	stub := &httpUpstreamStubFunc{
		doFunc: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
			capturedReq = req
			return upstreamResp, nil
		},
	}

	svc := &AntigravityGatewayService{
		httpUpstream: stub,
	}

	account := &Account{
		ID:       103,
		Name:     "upstream-node-gemini",
		Platform: PlatformAntigravity,
		Type:     AccountTypeUpstream,
		Credentials: map[string]any{
			"base_url": "http://154.36.173.146",
			"api_key":  "test-secret-key",
		},
	}

	body := []byte(`{"contents":[{"parts":[{"text":"hi"}]}]}`)
	res, err := svc.ForwardUpstreamGemini(context.Background(), c, account, "gemini-3-flash", "generateContent", false, body)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "gemini-3-flash", res.Model)
	require.Equal(t, 50, res.Usage.InputTokens)
	require.Equal(t, 100, res.Usage.OutputTokens)
	require.Equal(t, "http://154.36.173.146/v1beta/models/gemini-3-flash:generateContent", capturedReq.URL.String())
	require.Equal(t, "test-secret-key", capturedReq.Header.Get("x-goog-api-key"))
	require.Equal(t, "Bearer test-secret-key", capturedReq.Header.Get("Authorization"))
}

type httpUpstreamStubFunc struct {
	doFunc func(req *http.Request, proxyURL string, accountID int64, concurrency int) (*http.Response, error)
}

func (s *httpUpstreamStubFunc) Do(req *http.Request, proxyURL string, accountID int64, concurrency int) (*http.Response, error) {
	if s.doFunc != nil {
		return s.doFunc(req, proxyURL, accountID, concurrency)
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(nil))}, nil
}

func (s *httpUpstreamStubFunc) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return s.Do(req, proxyURL, accountID, concurrency)
}
