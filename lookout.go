package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log/slog"
	"net/mail"
	"os"
	"regexp"
	"strings"
)

// https://stackoverflow.com/a/57006724
const ISO8601 = "2006-01-02T15:04:05.999Z"

var webhookUrl string
var stampFilters []string

func main() {
	webhookUrl = os.Getenv("LOOKOUT_WEBHOOK")
	if webhookUrl == "" {
		slog.Error("invalid config option", slogTag("invalid_config"), slog.String("option", "LOOKOUT_WEBHOOK"))
		os.Exit(1)
	}
	stampsString := os.Getenv("LOOKOUT_STAMPS")
	stampFilters = strings.Split(stampsString, ",")
	if len(stampFilters) == 0 {
		slog.Warn("LOOKOUT_STAMPS is empty, all phabricator emails will be matched", slogTag("missing_stamps"))
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("failed to create watcher", slogTag("watcher_creation_failed"), slogError(err))
		os.Exit(1)
	}

	go watchLoop(watcher)

	err = watcher.Add("/mail")
	if err != nil {
		slog.Error("failed to watch /mail: is /mail volume mounted correctly?", slogTag("watcher_add_failed"), slogError(err))
		os.Exit(1)
	}

	<-make(chan struct{}) // Block forever
}

func watchLoop(w *fsnotify.Watcher) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				slog.Error("fsnotify watcher closed unexpectedly, exiting", slogTag("watcher_closed"), slog.String("channel", "w.Errors"))
				os.Exit(1)
			}
			slog.Error("fsnotify watcher error", slogTag("watcher_error"), slogError(err))
		case e, ok := <-w.Events:
			if !ok {
				slog.Error("fsnotify watcher closed unexpectedly, exiting", slogTag("watcher_closed"), slog.String("channel", "w.Events"))
				os.Exit(1)
			}
			logger := slog.With(slog.String("name", e.Name))
			logger.Debug("watcher got new event", slogTag("watcher_event"), slog.String("op", e.Op.String()))
			if e.Op.Has(fsnotify.Create) {
				go handleCreate(e, logger)
			}
		}
	}
}

var descriptionRegex = regexp.MustCompile(`(?s)(.+)\n\nTASK DETAIL\n {2}https://phabricator\.wikimedia\.org`)

// https://gerrit.wikimedia.org/r/plugins/gitiles/operations/puppet/+/703da43280b8ab00bd51602f83c38a108285cb25/modules/phabricator/data/fixed_settings.yaml#49
// https://secure.phabricator.com/uiexample/view/PHUIColorPalletteExample/
var priorityColours = map[string]int{
	"100": 0xda49be,
	"90":  0x8e44ad,
	"80":  0xc0392b,
	"50":  0xe67e22,
	"25":  0xf1c40f,
	"10":  0x3498db,
}

func handleCreate(e fsnotify.Event, logger *slog.Logger) {
	file, err := os.Open(e.Name)
	if err != nil {
		logger.Error("failed to open file", slogTag("open_failed"), slogError(err))
		return
	}
	parsedMail, err := ParseMail(file, logger)
	if err != nil {
		// error already logged in ParseMail
		return
	}
	stamps := parsedMail.Message.Header.Get("X-Phabricator-Stamps")
	subject := parsedMail.Message.Header.Get("Subject")
	if parsedMail.Message.Header.Get("X-Phabricator-Sent-This-Message") == "Yes" && strings.Contains(stamps, "application(Maniphest)") && MatchesStampFilter(stamps) {
		mailId := parsedMail.Message.Header.Get("X-Phabricator-Mail-ID")
		logger.Info("got maniphest email", slogTag("maniphest_email"), slog.String("mailId", mailId), slog.String("subject", subject), slog.String("stamps", stamps))
		logger = slog.With("mailId", mailId) // log with mailId instead of filename to save log space

		embed := Embed{
			Author: &EmbedAuthor{},
			Title:  subject[12:], // remove "[Maniphest] " from subject
			Fields: []*EmbedField{},
		}

		description := descriptionRegex.FindStringSubmatch(parsedMail.PlainBody)
		if len(description) < 2 {
			logger.Error("missing description", slogTag("missing_description"), slog.String("plainBody", parsedMail.PlainBody))
		} else {
			embed.Description = description[1]
		}

		if monogram := GetStamp("monogram", stamps, logger); monogram != "" {
			embed.URL = "https://phabricator.wikimedia.org/" + monogram
		}

		if actor := GetStamp("actor", stamps, logger); actor != "" {
			embed.Author.Name = actor
			embed.Author.URL = "https://phabricator.wikimedia.org/p/" + actor[1:] // remove initial @ from the actor for url
		}

		if status := GetStamp("status", stamps, logger); status != "" {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:   "Status",
				Value:  strings.ToUpper(status[:1]) + status[1:],
				Inline: true,
			})
		}

		if owner := GetStamp("task-owner", stamps, dummyLogger); owner != "" {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:   "Assignee",
				Value:  fmt.Sprintf("[%s](https://phabricator.wikimedia.org/p/%s)", owner, owner[1:]), // remove initial @ from the assignee for url
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:   "Assignee",
				Value:  "(none)",
				Inline: true,
			})
		}

		if tag := GetStamp("tag", stamps, dummyLogger); tag != "" {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:   "Tag",
				Value:  fmt.Sprintf("[%s](https://phabricator.wikimedia.org/p/%s)", tag, tag[1:]), // remove initial # from the tag for url
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:   "Tag",
				Value:  "(none)",
				Inline: true,
			})
		}

		if time, err := mail.ParseDate(parsedMail.Message.Header.Get("Date")); err != nil {
			slog.Error("failed to parse date header", slogTag("invalid_date"), slogError(err), slog.String("date", parsedMail.Message.Header.Get("Date")))
		} else {
			embed.Timestamp = time.Format(ISO8601)
		}

		if priority := GetStamp("task-priority", stamps, logger); priority != "" {
			embed.Color = priorityColours[priority]
		}

		executeWebhook(ExecuteWebhookPayload{
			Embeds: []Embed{embed},
		})
	}
}
