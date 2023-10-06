# lookout
Watches for Phabricator emails regarding certain tags and forwards them to a Discord webhook.

## config
All config is in environment variables.

| Key               | Description                                                                                                                                                                                                                                                                                                    | Type        |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|
| `LOOKOUT_WEBHOOK` | Discord webhook url                                                                                                                                                                                                                                                                                            | `string`    |
| `LOOKOUT_STAMPS`  | [Phabricator stamps](https://www.mediawiki.org/wiki/Phabricator/Help/Managing_mail#Filter_notifications) (`X-Phabricator-Stamps`, e.g. `tag(#ultarviolet)` or `subscriber(@Sportzpikachu)`) to watch for. Lookout will send a webhook if a phabricator email matches any of `LOOKOUT_STAMPS`. Comma seperated. | `[]string?` |
