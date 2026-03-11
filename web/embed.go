package web

import "embed"

//go:embed index.html assets/*
var StaticFS embed.FS
