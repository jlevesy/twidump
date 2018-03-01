package main

const templateContent = `
{{ range $tweet := .}}
From : {{$tweet.User.ScreenName}}  - Favorites : {{$tweet.FavoriteCount}} - Retweets : {{$tweet.RetweetCount}}
---------------------------------------------------
{{ if $tweet.Truncated }} Truncated {{ end }}
{{ if $tweet.Retweeted }} Retweeted {{ end }}
{{ if $tweet.ExtendedTweet }} Has Extended Tweet {{ end }}
{{ if $tweet.RetweetedStatus }} Has Retweeted Status {{ end }}

{{ if $tweet.Text }} {{ $tweet.Text }} {{ end }}
{{ if $tweet.FullText }} {{ $tweet.FullText }} {{ end }}
---------------------------------------------------
{{ else }}
No tweets for that user :(
{{ end }}
`
