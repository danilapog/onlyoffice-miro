package oauth

import (
	"context"
	"errors"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/storage/pg"
)

type OAuthable[T any] interface {
	Convert(T) (component.Authentication, error)
}

type oauthService[T any] struct {
	cipher         crypto.Cipher
	oauthClient    oauth.OAuthClient[T]
	oauthConverter OAuthable[T]
	storageService service.Storage[core.AuthCompositeKey, component.Authentication]
}

func NewOAuthService[T any](
	cipher crypto.Cipher,
	oauthClient oauth.OAuthClient[T],
	oauthConverter OAuthable[T],
	storageService service.Storage[core.AuthCompositeKey, component.Authentication],
) OAuthService[T] {
	return &oauthService[T]{
		cipher:         cipher,
		oauthClient:    oauthClient,
		oauthConverter: oauthConverter,
		storageService: storageService,
	}
}

func (s *oauthService[T]) encryptTokens(token component.Authentication) (string, string, error) {
	encAccess, err := s.cipher.Encrypt(token.AccessToken)
	if err != nil {
		return "", "", err
	}

	encRefresh, err := s.cipher.Encrypt(token.RefreshToken)
	if err != nil {
		return "", "", err
	}

	return encAccess, encRefresh, nil
}

func (s *oauthService[T]) decryptTokens(token component.Authentication) (string, string, error) {
	decAccess, err := s.cipher.Decrypt(token.AccessToken)
	if err != nil {
		return "", "", err
	}

	decRefresh, err := s.cipher.Decrypt(token.RefreshToken)
	if err != nil {
		return "", "", err
	}

	return decAccess, decRefresh, nil
}

func (s *oauthService[T]) createEncryptedAuth(token component.Authentication) (component.Authentication, error) {
	encAccess, encRefresh, err := s.encryptTokens(token)
	if err != nil {
		return component.Authentication{}, err
	}

	return component.Authentication{
		TokenType:    token.TokenType,
		AccessToken:  encAccess,
		RefreshToken: encRefresh,
		ExpiresAt:    token.ExpiresAt,
		Scope:        token.Scope,
	}, nil
}

func (s *oauthService[T]) createDecryptedAuth(token component.Authentication) (component.Authentication, error) {
	decAccess, decRefresh, err := s.decryptTokens(token)
	if err != nil {
		return component.Authentication{}, err
	}

	return component.Authentication{
		TokenType:    token.TokenType,
		AccessToken:  decAccess,
		RefreshToken: decRefresh,
		ExpiresAt:    token.ExpiresAt,
		Scope:        token.Scope,
	}, nil
}

func (s *oauthService[T]) Save(ctx context.Context, teamID, userID string, token component.Authentication) error {
	auth, err := s.createEncryptedAuth(token)
	if err != nil {
		return err
	}

	_, err = s.storageService.Insert(ctx, core.AuthCompositeKey{
		TeamID: teamID,
		UserID: userID,
	}, auth)
	return err
}

func (s *oauthService[T]) Find(ctx context.Context, teamID, userID string) (component.Authentication, error) {
	key := core.AuthCompositeKey{
		TeamID: teamID,
		UserID: userID,
	}

	storedAuth, err := s.storageService.Find(ctx, key)
	if err != nil {
		if errors.Is(err, pg.ErrNoRowsAffected) {
			return storedAuth, ErrTokenMissing
		}

		return storedAuth, err
	}

	if storedAuth.AccessToken == "" {
		return storedAuth, ErrTokenMissing
	}

	if time.Now().Unix() <= int64(storedAuth.ExpiresAt) {
		return s.createDecryptedAuth(storedAuth)
	}

	refreshToken, err := s.cipher.Decrypt(storedAuth.RefreshToken)
	if err != nil {
		return component.Authentication{}, err
	}

	token, err := s.oauthClient.Refresh(ctx, refreshToken)
	if err != nil {
		return component.Authentication{}, err
	}

	refreshedToken, err := s.oauthConverter.Convert(token)
	if err != nil {
		return component.Authentication{}, err
	}

	updatedAuth, err := s.createEncryptedAuth(refreshedToken)
	if err != nil {
		return component.Authentication{}, err
	}

	if _, err = s.storageService.Update(ctx, key, updatedAuth); err != nil {
		return component.Authentication{}, err
	}

	return refreshedToken, nil
}
