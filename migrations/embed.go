package migrations

import "embed"

// FS contiene las migraciones versionadas que consume goose.
//
//go:embed *.sql
var FS embed.FS
