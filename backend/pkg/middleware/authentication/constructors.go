package authentication

import (
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/labstack/echo/v4"
)

func NewHeaderAuthMiddleware(
	config *config.Config,
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
	headerName string,
) *AuthMiddleware {
	return NewAuthMiddleware(
		config,
		HeaderTokenExtractor(headerName),
		NoOpRefresher(),
		jwtService,
		translator,
		logger,
	)
}

func NewMiroAuthMiddleware(
	config *config.Config,
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
) *AuthMiddleware {
	middleware := NewAuthMiddleware(
		config,
		MiroSignatureExtractor(),
		NoOpRefresher(),
		jwtService,
		translator,
		logger,
	)

	extractor := func(c echo.Context) (string, error) {
		tokenString, err := MiroSignatureExtractor()(c)
		if err != nil {
			return "", err
		}

		token, err := middleware.ValidateToken(tokenString)
		if err != nil {
			return "", err
		}

		if err := middleware.SetAuthCookie(c, token.User, token.Team, int(token.ExpiresAt.Unix())); err != nil {
			return "", err
		}

		return tokenString, nil
	}

	middleware.extractor = extractor
	return middleware
}

func NewCookieOAuthMiddleware(
	config *config.Config,
	oauthService oauth.OAuthService[miro.AuthenticationResponse],
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
) *AuthMiddleware {
	middleware := NewAuthMiddleware(
		config,
		CookieTokenExtractor(config.Cookie.Name),
		nil,
		jwtService,
		translator,
		logger,
	)

	middleware.refresher = MiroOAuthTokenRefresher(middleware, oauthService)
	return middleware
}

func NewTokenAuthMiddleware(
	config *config.Config,
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
) *AuthMiddleware {
	return NewHeaderAuthMiddleware(config, jwtService, translator, logger, miroSignature)
}

func NewCookieAuthMiddleware(
	config *config.Config,
	oauthService oauth.OAuthService[miro.AuthenticationResponse],
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
) *AuthMiddleware {
	return NewCookieOAuthMiddleware(config, oauthService, jwtService, translator, logger)
}

func NewEditorAuthMiddleware(
	config *config.Config,
	oauthService oauth.OAuthService[miro.AuthenticationResponse],
	jwtService crypto.Signer,
	translator service.TranslationProvider,
	logger service.Logger,
) *AuthMiddleware {
	return NewCookieOAuthMiddleware(config, oauthService, jwtService, translator, logger)
}
