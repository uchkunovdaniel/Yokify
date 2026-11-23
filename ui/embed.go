package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:build
var files embed.FS

func FS() fs.FS {
	sub, _ := fs.Sub(files, "build")
	return sub
}
