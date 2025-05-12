package callback

type callbackQueryParams struct {
	UID string `json:"uid"`
	TID string `json:"tid"`
	BID string `json:"bid"`
	FID string `json:"fid"`
}

type callbackRequest struct {
	Status int    `json:"status"`
	Url    string `json:"url,omitempty"`
	Token  string `json:"token,omitempty"`
}

func (r *callbackRequest) Validate() error {
	if r.Token == "" {
		return ErrTokenRequired
	}

	return nil
}
