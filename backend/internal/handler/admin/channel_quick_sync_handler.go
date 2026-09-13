package admin

import (
	"errors"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// QuickSyncProbe handles probing upstream models with URL and Key.
// POST /api/v1/admin/channels/quick-sync/probe
func (h *ChannelHandler) QuickSyncProbe(c *gin.Context) {
	if h.quickSyncService == nil {
		response.ErrorFrom(c, infraerrors.InternalServer("QUICK_SYNC_UNAVAILABLE", "quick sync service not configured"))
		return
	}

	var params service.QuickSyncProbeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorFrom(c, infraerrors.BadRequest("VALIDATION_ERROR", err.Error()))
		return
	}

	result, err := h.quickSyncService.ProbeUpstream(c.Request.Context(), params)
	if err != nil {
		var appErr *infraerrors.ApplicationError
		if errors.As(err, &appErr) {
			response.ErrorFrom(c, appErr)
		} else {
			response.ErrorFrom(c, infraerrors.BadRequest("PROBE_FAILED", err.Error()))
		}
		return
	}

	response.Success(c, result)
}

// QuickSyncCommit handles committing the quick sync channel onboarding.
// POST /api/v1/admin/channels/quick-sync/commit
func (h *ChannelHandler) QuickSyncCommit(c *gin.Context) {
	if h.quickSyncService == nil {
		response.ErrorFrom(c, infraerrors.InternalServer("QUICK_SYNC_UNAVAILABLE", "quick sync service not configured"))
		return
	}

	var params service.QuickSyncCommitParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorFrom(c, infraerrors.BadRequest("VALIDATION_ERROR", err.Error()))
		return
	}

	result, err := h.quickSyncService.CommitQuickSync(c.Request.Context(), params)
	if err != nil {
		var appErr *infraerrors.ApplicationError
		if errors.As(err, &appErr) {
			response.ErrorFrom(c, appErr)
		} else {
			response.ErrorFrom(c, infraerrors.BadRequest("COMMIT_FAILED", err.Error()))
		}
		return
	}

	response.Success(c, result)
}
