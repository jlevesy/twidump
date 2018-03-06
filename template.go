package main

const templateContent = `
{{- range $tweet := . -}}
{{$tweet.CreatedAt}} From : {{$tweet.User.ScreenName}}  - Favorites : {{$tweet.FavoriteCount}} - Retweets : {{$tweet.RetweetCount}}
---------------------------------------------------
{{- if $tweet.Retweeted }}
RT @{{ $tweet.RetweetedStatus.User.ScreenName}}: {{$tweet.RetweetedStatus.FullText}}
{{- else }}
{{ with $tweet.QuotedStatus -}}
Quote from @{{ .User.ScreenName}}: {{.FullText}}
{{ end }}
{{ $tweet.FullText }}
{{ with $tweet.ExtendedEntities }}
Media :
{{- range $medium := .Media }}
Type: {{$medium.Type}} - Display URL {{$medium.DisplayURL}} - URL : {{$medium.MediaURL}}
{{- end -}}
{{- end -}}
{{ end }}
---------------------------------------------------
{{ else }}
No tweets for that user :(
{{ end -}}
`
