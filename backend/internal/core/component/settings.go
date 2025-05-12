package component

type Settings struct {
	Address string `json:"address"`
	Header  string `json:"header"`
	Secret  string `json:"secret"`
	Demo    Demo   `json:"demo,omitempty"`
}
