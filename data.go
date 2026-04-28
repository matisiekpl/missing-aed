package backend

import "embed"

//go:embed competition/*.json
var CompetitionFiles embed.FS
