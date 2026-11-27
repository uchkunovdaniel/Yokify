// Package ui embeds the ui
package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:build
var files embed.FS

// FS returns the embedded UI files.
func FS() fs.FS {
	sub, _ := fs.Sub(files, "build")
	return sub
}
