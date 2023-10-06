package main

import (
	"encoding/base64"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

type Mail struct {
	Message   *mail.Message
	PlainBody string
}

func ParseMail(r io.Reader, logger *slog.Logger) (parsedMail Mail, err error) {
	message, err := mail.ReadMessage(r)
	if err != nil {
		logger.Error("failed to read message", slogTag("message_read_failed"), slogError(err))
		return
	}
	parsedMail.Message = message

	mediaType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
	if err != nil {
		logger.Error("failed to parse Content-Type", slogTag("media_type_invalid"), slogError(err))
		return
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(message.Body, params["boundary"])
		for {
			var p *multipart.Part
			p, err = mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error("failed to get part", slogTag("part_invalid"), slogError(err))
				return
			}
			var mediaType string
			mediaType, _, err = mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				logger.Error("failed to parse Content-Type", slogTag("media_type_invalid"), slogError(err))
				return
			}
			if mediaType == "text/plain" {
				var reader io.Reader
				if p.Header.Get("Content-Transfer-Encoding") == "base64" {
					reader = base64.NewDecoder(base64.StdEncoding, p)
				} else {
					reader = p
				}

				var bytes []byte
				bytes, err = io.ReadAll(reader)
				if err != nil {
					slog.Error("failed to read part", slogTag("part_read_failed"), slogError(err))
				}
				parsedMail.PlainBody = string(bytes)
				break
			}
		}
	}

	return
}
