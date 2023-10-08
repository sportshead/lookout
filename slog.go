package main

import (
	"io"
	"log/slog"
)

func slogTag(tag string) slog.Attr {
	return slog.String("tag", tag)
}

func slogError(err error) slog.Attr {
	return slog.Any("error", err)
}

var dummyLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
