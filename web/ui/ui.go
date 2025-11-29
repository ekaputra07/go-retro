package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var UiFS embed.FS

func EmbeddedUiFS() http.FileSystem {
	staticFiles, err := fs.Sub(UiFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(staticFiles)
}
