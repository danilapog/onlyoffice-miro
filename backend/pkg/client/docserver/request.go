package docserver

type GetServerVersionRequest struct {
	C     string `json:"c"`
	Token string `json:"token,omitempty"`
}
