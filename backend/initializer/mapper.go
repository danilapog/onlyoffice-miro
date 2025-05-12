package initializer

import (
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
)

type authenticationMapper struct{}

func NewAuthenticationMapper() *authenticationMapper {
	return &authenticationMapper{}
}

func (m *authenticationMapper) Convert(token miro.AuthenticationResponse) (component.Authentication, error) {
	expiresAt := time.Now().Add(time.Second * time.Duration(token.ExpiresIn)).Unix()
	return component.Authentication{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    int(expiresAt),
		Scope:        token.Scope,
	}, nil
}
