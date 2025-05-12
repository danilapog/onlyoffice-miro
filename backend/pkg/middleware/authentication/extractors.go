package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/labstack/echo/v4"
)

func HeaderTokenExtractor(headerName string) TokenExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Request().Header.Get(headerName)
		if token == "" {
			return "", echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Missing %s header", headerName))
		}
		return token, nil
	}
}

func CookieTokenExtractor(cookieName string) TokenExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			signature := c.Request().Header.Get(miroSignature)
			if signature != "" {
				return signature, nil
			}
			return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing authentication cookie")
		}
		if cookie.Value == "" {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "Empty authentication cookie")
		}
		return cookie.Value, nil
	}
}

func MiroSignatureExtractor() TokenExtractor {
	return HeaderTokenExtractor(miroSignature)
}

func NoOpRefresher() TokenRefresher {
	return func(c echo.Context, token *TokenClaims) error {
		return nil
	}
}

func MiroOAuthTokenRefresher(middleware *AuthMiddleware, oauthService oauth.OAuthService[miro.AuthenticationResponse]) TokenRefresher {
	return func(c echo.Context, token *TokenClaims) error {
		if token.ExpiresAt.Before(time.Now()) || time.Until(token.ExpiresAt.Time) < time.Hour {
			_, err := oauthService.Find(c.Request().Context(), token.Team, token.User)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to refresh token")
			}

			expiresAt := int(time.Now().Add(24 * time.Hour).Unix())
			if err := middleware.SetAuthCookie(c, token.User, token.Team, expiresAt); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set auth cookie")
			}
		}

		return nil
	}
}
