package docserver

type ServerVersionResponse struct {
	Error   int    `json:"error"`
	Version string `json:"version"`
}
