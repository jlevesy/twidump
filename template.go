package main

const templateContent = `
{{- range $tweet := . -}}
From : {{$tweet.User.ScreenName}}  - Favorites : {{$tweet.FavoriteCount}} - Retweets : {{$tweet.RetweetCount}}
---------------------------------------------------
{{- if $tweet.Retweeted }}
RT @{{ $tweet.RetweetedStatus.User.ScreenName}}: {{$tweet.RetweetedStatus.FullText}}
{{- else }}
{{ with $tweet.QuotedStatus -}}
Quote from @{{ .User.ScreenName}}: {{.FullText}}
{{ end }}
{{ $tweet.FullText }}
{{ end }}
---------------------------------------------------
{{ else }}
 No tweets for that user :(
{{ end -}}
`
