package main

const templateContent = `
{{ range $tweet := .Tweets}}
From : {{$tweet.User.ScreenName}}  - Favorites : {{$tweet.FavoriteCount}} - Retweets : {{$tweet.RetweetCount}}
---------------------------------------------------
{{$tweet.Text}}
---------------------------------------------------
{{ end }}
`