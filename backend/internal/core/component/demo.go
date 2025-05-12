package component

import "time"

type Demo struct {
	TeamID  string     `json:"team_id"`
	Enabled bool       `json:"enabled"`
	Started *time.Time `json:"started,omitempty"`
}
