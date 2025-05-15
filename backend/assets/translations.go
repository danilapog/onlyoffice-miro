package assets

import "embed"

//go:embed translations/*.json
var Translations embed.FS
