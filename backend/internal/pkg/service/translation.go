package service

import (
	"context"
)

type TranslationProvider interface {
	Translate(ctx context.Context, lang, id string) string
}
