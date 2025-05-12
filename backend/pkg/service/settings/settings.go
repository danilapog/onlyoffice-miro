package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/storage/pg"
)

const defaultCacheExpiration = 5 * time.Minute

type settingsService struct {
	cipher         crypto.Cipher
	storageService service.Storage[core.SettingsCompositeKey, component.Settings]
	cache          service.Cache
	logger         service.Logger
}

func NewSettingsService(
	cipher crypto.Cipher,
	storageService service.Storage[core.SettingsCompositeKey, component.Settings],
	cache service.Cache,
	logger service.Logger,
) SettingsService {
	return &settingsService{
		cipher:         cipher,
		storageService: storageService,
		cache:          cache,
		logger:         logger,
	}
}

func (s *settingsService) buildCacheKey(teamID, boardID string) string {
	return fmt.Sprintf("settings:%s:%s", teamID, boardID)
}

func (s *settingsService) createCompositeKey(teamID, boardID string) core.SettingsCompositeKey {
	return core.SettingsCompositeKey{
		TeamID:  teamID,
		BoardID: boardID,
	}
}

func (s *settingsService) logEvent(ctx context.Context, level config.Level, message string, teamID, boardID string, err error) {
	fields := service.Fields{
		"team_id":  teamID,
		"board_id": boardID,
	}

	if err != nil {
		fields["error"] = err.Error()
	}

	switch level {
	case config.Debug:
		s.logger.Debug(ctx, message, fields)
	case config.Warn:
		s.logger.Warn(ctx, message, fields)
	case config.Error:
		s.logger.Error(ctx, message, fields)
	}
}

func (s *settingsService) encryptSecret(secret string) (string, error) {
	if secret == "" {
		return "", nil
	}

	return s.cipher.Encrypt(secret)
}

func (s *settingsService) decryptSecret(encryptedSecret string) (string, error) {
	if encryptedSecret == "" {
		return "", nil
	}

	return s.cipher.Decrypt(encryptedSecret)
}

func (s *settingsService) createDemoSettings(teamID string, existingStarted *time.Time) component.Demo {
	demoSettings := component.Demo{
		TeamID:  teamID,
		Enabled: true,
	}

	if existingStarted != nil {
		demoSettings.Started = existingStarted
	} else {
		now := time.Now()
		demoSettings.Started = &now
	}

	return demoSettings
}

func (s *settingsService) invalidateCache(ctx context.Context, teamID, boardID string) {
	cacheKey := s.buildCacheKey(teamID, boardID)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logEvent(ctx, config.Warn, "Failed to invalidate settings cache", teamID, boardID, err)
	}
}

func (s *settingsService) cacheSettings(ctx context.Context, teamID, boardID string, settings component.Settings) {
	cacheKey := s.buildCacheKey(teamID, boardID)
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		s.logEvent(ctx, config.Warn, "Failed to marshal settings for caching", teamID, boardID, err)
		return
	}

	if err := s.cache.Set(ctx, cacheKey, settingsJson, defaultCacheExpiration); err != nil {
		s.logEvent(ctx, config.Warn, "Failed to cache settings", teamID, boardID, err)
	}
}

func (s *settingsService) Save(ctx context.Context, teamID, boardID string, opts ...Option) error {
	settings := &SaveOptions{}
	for _, opt := range opts {
		opt(settings)
	}

	if err := settings.Validate(); err != nil {
		return err
	}

	compositeKey := s.createCompositeKey(teamID, boardID)
	existingSettings, err := s.storageService.Find(ctx, compositeKey)
	if err != nil && !errors.Is(err, pg.ErrNoRowsAffected) {
		return err
	}

	newSettings, err := s.buildNewSettings(teamID, settings, existingSettings)
	if err != nil {
		return err
	}

	if _, err := s.storageService.Insert(ctx, compositeKey, newSettings); err != nil {
		return err
	}

	s.invalidateCache(ctx, teamID, boardID)

	return nil
}

func (s *settingsService) buildNewSettings(teamID string, opts *SaveOptions, existingSettings component.Settings) (component.Settings, error) {
	var newSettings component.Settings

	if opts.Demo {
		newSettings.Demo = s.createDemoSettings(teamID, existingSettings.Demo.Started)

		if opts.Address != "" || opts.Header != "" || opts.Secret != "" {
			encSecret, err := s.encryptSecret(opts.Secret)
			if err != nil {
				return newSettings, err
			}

			newSettings.Address = opts.Address
			newSettings.Header = opts.Header
			newSettings.Secret = encSecret
		}

		return newSettings, nil
	}

	encSecret, err := s.encryptSecret(opts.Secret)
	if err != nil {
		return newSettings, err
	}

	newSettings = component.Settings{
		Address: opts.Address,
		Header:  opts.Header,
		Secret:  encSecret,
	}

	if existingSettings.Demo.Enabled {
		newSettings.Demo = existingSettings.Demo
	}

	return newSettings, nil
}

func (s *settingsService) Find(ctx context.Context, teamID, boardID string) (component.Settings, error) {
	settings, found, err := s.getFromCache(ctx, teamID, boardID)
	if err != nil {
		s.logEvent(ctx, config.Warn, "Error processing cached settings", teamID, boardID, err)
	}

	if found {
		return settings, nil
	}

	return s.getFromStorage(ctx, teamID, boardID)
}

func (s *settingsService) getFromCache(ctx context.Context, teamID, boardID string) (component.Settings, bool, error) {
	cacheKey := s.buildCacheKey(teamID, boardID)
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return component.Settings{}, false, err
	}

	if cachedData == nil {
		return component.Settings{}, false, nil
	}

	var settings component.Settings
	if err := json.Unmarshal(cachedData, &settings); err != nil {
		return component.Settings{}, false, err
	}

	if settings.Demo.Enabled && (settings.Address == "" || settings.Header == "" || settings.Secret == "") {
		s.logEvent(ctx, config.Debug, "Demo settings retrieved from cache", teamID, boardID, nil)
		return settings, true, nil
	}

	if settings.Secret != "" {
		decSecret, err := s.decryptSecret(settings.Secret)
		if err != nil {
			return settings, false, err
		}
		settings.Secret = decSecret
	}

	s.logEvent(ctx, config.Debug, "Settings retrieved from cache", teamID, boardID, nil)
	return settings, true, nil
}

func (s *settingsService) getFromStorage(ctx context.Context, teamID, boardID string) (component.Settings, error) {
	compositeKey := s.createCompositeKey(teamID, boardID)
	settings, err := s.storageService.Find(ctx, compositeKey)

	if err != nil {
		if errors.Is(err, pg.ErrNoRowsAffected) {
			return component.Settings{}, nil
		}
		return component.Settings{}, err
	}

	s.cacheSettings(ctx, teamID, boardID, settings)

	if settings.Demo.Enabled && (settings.Address == "" || settings.Header == "" || settings.Secret == "") {
		return settings, nil
	}

	if settings.Secret != "" {
		decSecret, err := s.decryptSecret(settings.Secret)
		if err != nil {
			return settings, err
		}
		settings.Secret = decSecret
	}

	return settings, nil
}
