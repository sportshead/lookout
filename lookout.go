package main

import (
	"log/slog"
	"os"
	"strings"
)

var webhookUrl string
var stamps []string

func main() {
	webhookUrl = os.Getenv("LOOKOUT_WEBHOOK")
	if webhookUrl == "" {
		slog.Error("invalid config option", slogTag("invalid_config"), slog.String("option", "LOOKOUT_WEBHOOK"))
	}
	stampsString := os.Getenv("LOOKOUT_STAMPS")
	stamps = strings.Split(stampsString, ",")
	if len(stamps) == 0 {
		slog.Warn("LOOKOUT_STAMPS is empty, all phabricator emails will be matched", slogTag("missing_stamps"))
	}
}
