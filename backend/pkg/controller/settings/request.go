package settings

type settingsRequest struct {
	BoardID string `json:"board_id"`
	Address string `json:"address"`
	Header  string `json:"header"`
	Secret  string `json:"secret"`
	Demo    bool   `json:"demo"`
}

func (r *settingsRequest) Validate() error {
	if r.BoardID == "" {
		return ErrBoardIdRequired
	}

	return nil
}

type persistSettingsRequest struct {
	BoardID string `json:"board_id"`
	Address string `json:"address"`
	Header  string `json:"header"`
	Secret  string `json:"secret"`
	Demo    bool   `json:"demo"`
}

func (r *persistSettingsRequest) Validate() error {
	if r.BoardID == "" {
		return ErrBoardIdRequired
	}

	return nil
}
