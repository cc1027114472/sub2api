package service

import (
	"context"
	"log/slog"
	"time"
)

func (s *ContentModerationService) applyFlaggedAccountSideEffects(ctx context.Context, cfg *ContentModerationConfig, log *ContentModerationLog) bool {
	if s == nil || cfg == nil || log == nil || !log.Flagged || log.UserID == nil || *log.UserID <= 0 {
		return false
	}
	count := 1
	if s.repo != nil && cfg.ViolationWindowHours > 0 {
		since := time.Now().Add(-time.Duration(cfg.ViolationWindowHours) * time.Hour)
		if n, err := s.repo.CountFlaggedByUserSince(ctx, *log.UserID, since, cfg.CyberPolicyExcludeFromBanCount); err == nil {
			count = n + 1
		}
	}
	log.ViolationCount = count
	autoBanJustApplied := false
	if cfg.AutoBanEnabled && cfg.BanThreshold > 0 && count >= cfg.BanThreshold && s.userRepo != nil {
		user, err := s.userRepo.GetByID(ctx, *log.UserID)
		if err != nil {
			slog.Warn("content_moderation.ban_get_user_failed", "user_id", *log.UserID, "error", err)
			return false
		}
		if user.IsAdmin() {
			slog.Warn("content_moderation.autoban_skipped_admin", "user_id", *log.UserID, "role", user.Role, "count", count, "threshold", cfg.BanThreshold)
			s.disableTriggeringAPIKeyForAdmin(ctx, log)
			return false
		}
		if user.Status != StatusDisabled {
			user.Status = StatusDisabled
			if err := s.userRepo.Update(ctx, user, UserUpdateFields{Status: true}); err != nil {
				slog.Warn("content_moderation.ban_update_user_failed", "user_id", *log.UserID, "error", err)
				return false
			}
			if s.authCacheInvalidator != nil {
				s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, *log.UserID)
			}
			autoBanJustApplied = true
		}
		log.AutoBanned = true
	}
	return autoBanJustApplied
}

func (s *ContentModerationService) disableTriggeringAPIKeyForAdmin(ctx context.Context, log *ContentModerationLog) {
	if s == nil || log == nil || log.APIKeyID == nil || *log.APIKeyID <= 0 || s.apiKeyRepo == nil {
		return
	}
	apiKey, err := s.apiKeyRepo.GetByID(ctx, *log.APIKeyID)
	if err != nil {
		slog.Warn("content_moderation.admin_disable_api_key_get_failed", "api_key_id", *log.APIKeyID, "error", err)
		return
	}
	if apiKey.Status == StatusAPIKeyDisabled {
		return
	}
	apiKey.Status = StatusAPIKeyDisabled
	if err := s.apiKeyRepo.Update(ctx, apiKey, APIKeyUpdateFields{Status: true}); err != nil {
		slog.Warn("content_moderation.admin_disable_api_key_update_failed", "api_key_id", *log.APIKeyID, "error", err)
		return
	}
	if s.authCacheInvalidator != nil && apiKey.Key != "" {
		s.authCacheInvalidator.InvalidateAuthCacheByKey(ctx, apiKey.Key)
	}
	slog.Warn("content_moderation.admin_api_key_disabled", "user_id", contentModerationEmailUserID(log), "api_key_id", *log.APIKeyID)
}

