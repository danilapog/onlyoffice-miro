package file

import "github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"

type convertResponse struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type boardAuthenticationResponse struct {
	BoardID        string                    `json:"board_id"`
	Authentication *component.Authentication `json:"authentication"`
}
