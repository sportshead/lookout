package main

import (
	"embed"
	"log/slog"
	"testing"
)

//go:embed sample_emails/*
var f embed.FS

func TestParseMail_simple(t *testing.T) {
	logger := slog.With()
	simpleEmail, err := f.Open("sample_emails/simple.eml")
	expectedBody, err := f.ReadFile("sample_emails/simple.eml.txt")
	if err != nil {
		t.Fatal("reading simple.eml failed", err)
	}
	parsedMail, err := ParseMail(simpleEmail, logger)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if parsedMail.PlainBody == "" {
		t.Error("plain body is blank")
	}
	if parsedMail.PlainBody != string(expectedBody) {
		t.Errorf("expected in sample_emails/simple.eml.txt, got %s", parsedMail.PlainBody)
	}
}
func TestParseMail_comment(t *testing.T) {
	logger := slog.With()
	simpleEmail, err := f.Open("sample_emails/comment.eml")
	expectedBody, err := f.ReadFile("sample_emails/comment.eml.txt")
	if err != nil {
		t.Fatal("reading comment.eml failed", err)
	}
	parsedMail, err := ParseMail(simpleEmail, logger)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if parsedMail.PlainBody == "" {
		t.Error("plain body is blank")
	}
	if parsedMail.PlainBody != string(expectedBody) {
		t.Errorf("expected in sample_emails/comment.eml.txt, got %s", parsedMail.PlainBody)
	}
}
