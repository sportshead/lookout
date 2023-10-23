# lookout
[![Go Report Card](https://goreportcard.com/badge/gitlab.wikimedia.org/sportz/lookout)](https://goreportcard.com/report/gitlab.wikimedia.org/sportz/lookout)
[![License](https://img.shields.io/gitlab/license/sportz%2Flookout?gitlab_url=https%3A%2F%2Fgitlab.wikimedia.org)](LICENSE)
[![Gitlab pipeline status](https://img.shields.io/gitlab/pipeline-status/sportz%2Flookout?gitlab_url=https%3A%2F%2Fgitlab.wikimedia.org&branch=master&logo=gitlab)](https://gitlab.wikimedia.org/sportz/lookout/-/pipelines)

Watches for Phabricator emails regarding certain tags and forwards them to a Discord webhook.

## usage
All config is in environment variables. Make sure to bind mount the directory which new mail will appear in to `/mail` on the container.

| Key               | Description                                                                                                                                                                                                                                                                                                    | Type        |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|
| `LOOKOUT_WEBHOOK` | Discord webhook url                                                                                                                                                                                                                                                                                            | `string`    |
| `LOOKOUT_STAMPS`  | [Phabricator stamps](https://www.mediawiki.org/wiki/Phabricator/Help/Managing_mail#Filter_notifications) (`X-Phabricator-Stamps`, e.g. `tag(#ultarviolet)` or `subscriber(@Sportzpikachu)`) to watch for. Lookout will send a webhook if a phabricator email matches any of `LOOKOUT_STAMPS`. Comma seperated. | `[]string?` |

## planned features
- [ ] reply to message based on `Thread-Topic` header
  - [ ] persistent storage (sqlite?) linking `Thread-Topic`s to latest Discord message id

## things i learned in this project
- go packages
  - file watching with [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)
  - logging with [log/slog](https://pkg.go.dev/log/slog)
  - go unit testing with [testing](https://pkg.go.dev/testing)
  - mail parsing with [net/mail](https://pkg.go.dev/net/mail)
- gitlab ci/cd
  - setting up a personal runner
  - docker in docker
  - building docker in gitlab ci
