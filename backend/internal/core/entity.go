package core

type Entity string

type AuthCompositeKey struct {
	TeamID string
	UserID string
}

type SettingsCompositeKey struct {
	TeamID  string
	BoardID string
}
