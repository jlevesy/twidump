package main

const templateContent = `
{{- range $tweet := . -}}
{{$tweet.CreatedAt}} From : {{$tweet.User.ScreenName}}  - Favorites : {{$tweet.FavoriteCount}} - Retweets : {{$tweet.RetweetCount}}
---------------------------------------------------
{{ with $tweet.QuotedStatus -}}
Quote from @{{ .User.ScreenName}}: {{.FullText}}

{{ end -}}
{{with $tweet.RetweetedStatus -}}
RT @{{.User.ScreenName}}: {{ .FullText }}
{{- else -}}
{{ $tweet.FullText }}
{{- end -}}
{{ with $tweet.ExtendedEntities -}}
Media :
{{- range $medium := .Media }}
Type: {{$medium.Type}} - Display URL {{$medium.DisplayURL}} - URL : {{$medium.MediaURL}}
{{- end -}}
{{ end }}
---------------------------------------------------
{{ else }}
No tweets for that user :(
{{ end -}}
`
