package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

func GetStamp(stampKey string, stamps string, logger *slog.Logger) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?U)%s\((.+)\)`, stampKey))
	stamp := re.FindStringSubmatch(stamps)
	if len(stamp) < 2 {
		logger.Warn("missing stamp", slogTag("missing_stamp"), slog.String("stamp", stampKey))
		return ""
	}
	return stamp[1]
}

func MatchesStampFilter(stamps string) bool {
	for _, stampFilter := range stampFilters {
		if strings.Contains(stamps, stampFilter) {
			return true
		}
	}
	return false
}
