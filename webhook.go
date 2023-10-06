package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Snowflake string

type AllowedMentions struct {
	Parse       []AllowedMentionType `json:"parse,omitempty"`
	Roles       []string             `json:"roles,omitempty"`
	Users       []string             `json:"users,omitempty"`
	RepliedUser bool                 `json:"replied_user,omitempty"`
}

type AllowedMentionType string

const (
	AllowedMentionTypeRole     AllowedMentionType = "roles"
	AllowedMentionTypeUser     AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)

var DefaultAllowedMentions = &AllowedMentions{
	Parse: []AllowedMentionType{},
}

type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []*EmbedField   `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type ExecuteWebhookPayload struct {
	Content         string           `json:"content,omitempty"`
	Username        string           `json:"username,omitempty"`
	AvatarURL       string           `json:"avatar_url,omitempty"`
	TTS             bool             `json:"tts,omitempty"`
	Embeds          []Embed          `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int              `json:"flags,omitempty"`
	ThreadName      string           `json:"thread_name,omitempty"`
}

type ExecuteWebhookResult struct {
	ID Snowflake `json:"id"`
}

var client = http.Client{}

// Executes the webhook with the provided body. Returns the message snowflake
func executeWebhook(payload ExecuteWebhookPayload) Snowflake {
	url := webhookUrl + "?wait=true"
	logger := slog.With(slog.String("func", "executeWebhook"), slog.String("url", url))

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		logger.Error("json marshal failed", slogTag("failed_marshal"), slogError(err))
		return ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
	if err != nil {
		logger.Error("failed to create request", slogTag("failed_create_request"), slogError(err))
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("failed to send request", slogTag("failed_send_request"), slogError(err))
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logger.Error("response code not ok", slogTag("bad_response_code"), slog.Int("code", res.StatusCode), slog.String("message", res.Status))
		return ""
	}

	bodyJson, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("failed to read body", slogTag("failed_read_body"), slogError(err))
		return ""
	}

	body := &ExecuteWebhookResult{}
	err = json.Unmarshal(bodyJson, body)
	if err != nil {
		logger.Error("failed to unmarshal body", slogTag("failed_unmarshal_body"), slogError(err), slog.String("json", string(bodyJson)))
		return ""
	}

	slog.Info("webhook execution successful", slogTag("webhook_success"), slog.String("message_id", string(body.ID)))
	return body.ID
}
