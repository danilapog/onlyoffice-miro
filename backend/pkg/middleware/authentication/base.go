package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// TODO: Refactor, move to commons

type TokenClaims struct {
	User string `json:"user"`
	Team string `json:"team"`
	jwt.RegisteredClaims
}

type TokenExtractor func(c echo.Context) (string, error)

type TokenRefresher func(c echo.Context, token *TokenClaims) error

type AuthMiddleware struct {
	config     *config.Config
	jwtService crypto.Signer
	extractor  TokenExtractor
	refresher  TokenRefresher
	logger     service.Logger
}

func NewAuthMiddleware(
	config *config.Config,
	jwtService crypto.Signer,
	extractor TokenExtractor,
	refresher TokenRefresher,
	logger service.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		config:     config,
		jwtService: jwtService,
		extractor:  extractor,
		refresher:  refresher,
		logger:     logger,
	}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := m.extractor(c)
		lang := "en"
		if langParam := c.QueryParam("lang"); langParam != "" {
			lang = langParam
		}

		if err != nil {
			return c.Render(http.StatusOK, "unauthorized", map[string]string{
				"language":           lang,
				"authorizationError": "Missing authentication",
			})
		}

		m.logger.Info(c.Request().Context(), "authenticating request",
			service.Fields{
				"token": tokenString,
			})

		token, err := m.ValidateToken(tokenString)
		if err != nil {
			return c.Render(http.StatusOK, "unauthorized", map[string]string{
				"language":           lang,
				"authorizationError": "Missing authentication",
			})
		}

		if m.refresher != nil {
			if err := m.refresher(c, token); err != nil {
				return c.Render(http.StatusOK, "unauthorized", map[string]string{
					"language":           lang,
					"authorizationError": "Missing authentication",
				})
			}
		}

		m.logger.Info(c.Request().Context(), "authenticated request",
			service.Fields{
				"token": tokenString,
			})

		c.Set(common.ContextKeyUser, token)
		return next(c)
	}
}

func (m *AuthMiddleware) ValidateToken(tokenString string) (*TokenClaims, error) {
	var token TokenClaims
	if err := m.jwtService.ValidateTarget(tokenString, []byte(m.config.OAuth.ClientSecret), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (m *AuthMiddleware) CreateAuthToken(uid, tid string, expiresAt int) (string, error) {
	claims := &TokenClaims{
		User: uid,
		Team: tid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(int64(expiresAt), 0)),
		},
	}

	return m.jwtService.Create(claims, []byte(m.config.OAuth.ClientSecret))
}

func (m *AuthMiddleware) CreateAuthCookie(tokenString string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     m.config.Cookie.Name,
		Value:    tokenString,
		MaxAge:   m.config.Cookie.MaxAge,
		Path:     m.config.Cookie.Path,
		Domain:   m.config.Server.Domain,
		Secure:   m.config.Cookie.Secure,
		HttpOnly: m.config.Cookie.HttpOnly,
		SameSite: m.config.Cookie.GetSameSite(),
	}

	return cookie
}

func (m *AuthMiddleware) LogCookieInfo(cookie *http.Cookie, c echo.Context) {
	if cookie == nil {
		return
	}

	ctx := context.Background()
	if c != nil {
		ctx = c.Request().Context()
	}

	fields := service.Fields{
		"name":     cookie.Name,
		"domain":   cookie.Domain,
		"path":     cookie.Path,
		"secure":   cookie.Secure,
		"httpOnly": cookie.HttpOnly,
		"sameSite": m.config.Cookie.SameSite,
	}

	if c != nil {
		fields["method"] = c.Request().Method
		fields["path"] = c.Request().URL.Path
		fields["host"] = c.Request().Host
		fields["remote"] = c.Request().RemoteAddr
	}

	m.logger.Info(ctx, "cookie operation", fields)
}

func (m *AuthMiddleware) SetAuthCookie(c echo.Context, uid, tid string, expiresAt int) error {
	tokenString, err := m.CreateAuthToken(uid, tid, expiresAt)
	if err != nil {
		return err
	}

	cookie := m.CreateAuthCookie(tokenString)
	c.SetCookie(cookie)
	m.LogCookieInfo(cookie, c)

	return nil
}

func (m *AuthMiddleware) ClearAuthCookie(c echo.Context) {
	cookie := m.CreateAuthCookie("")
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	m.LogCookieInfo(cookie, nil)
}

func (m *AuthMiddleware) GetCookieExpiration(c echo.Context) error {
	tokenString, err := m.extractor(c)
	lang := "en"
	if langParam := c.QueryParam("lang"); langParam != "" {
		lang = langParam
	}

	if err != nil {
		return c.Render(http.StatusOK, "unauthorized", map[string]string{
			"language":           lang,
			"authorizationError": "Missing authentication",
		})
	}

	token, err := m.ValidateToken(tokenString)
	if err != nil {
		return c.Render(http.StatusOK, "unauthorized", map[string]string{
			"language":           lang,
			"authorizationError": "Invalid token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"expires_at": token.ExpiresAt.Unix(),
	})
}
