package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	concurrency         = 5
	credentialsFilePath = "credentials.json"

	defaultAccountsFilename = "accounts.csv"

	outputFilenameTemplate   = "tweedump-result-%s.txt"
	outputFilenameDateFormat = "2-01-2006-15-04-05"

	commaSeparator = ';'
)

type credentials struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
}

type account struct {
	ScreenName string
	SinceID    int64
}

type config struct {
	AccountsFile string
	OutputFile   string
}

func getConfig(args []string) *config {
	result := config{AccountsFile: defaultAccountsFilename}

	if len(args) > 1 {
		result.AccountsFile = args[1]
	}

	result.OutputFile = fmt.Sprintf(
		outputFilenameTemplate,
		time.Now().Format(outputFilenameDateFormat),
	)

	return &result
}

func setupClient(credentialsFilePath string) (*twitter.Client, error) {
	creds := credentials{}

	f, err := os.Open(credentialsFilePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&creds); err != nil {
		return nil, err
	}

	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient), nil
}

func fetchTimeline(client *twitter.Client, in <-chan *account, out chan<- []twitter.Tweet, wg *sync.WaitGroup) {
	for {
		account, ok := <-in

		if !ok {
			wg.Done()
			return
		}

		userTimelineParams := &twitter.UserTimelineParams{
			ScreenName: account.ScreenName,
			SinceID:    account.SinceID,
			TweetMode:  "extended",
		}

		tweets, _, err := client.Timelines.UserTimeline(userTimelineParams)

		if err != nil {
			log.Printf("Failed to collect tweets for user %s, reason is : %s", account.ScreenName, err)
		}

		out <- tweets
	}
}

func loadAccounts(accountFilePath string, in chan<- *account) {
	file, err := os.Open(accountFilePath)

	if err != nil {
		log.Fatal("Failed to open input file, reason is :", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = commaSeparator

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		sinceID, err := strconv.ParseInt(record[1], 10, 64)

		if err != nil {
			log.Printf("Warning, failed to parse sinceID for user %s", record[0])
			sinceID = 0
		}

		in <- &account{record[0], sinceID}
	}
}

func dumpTimeline(outputFilePath string, out <-chan []twitter.Tweet, done chan<- struct{}) {
	t := template.Must(template.New("timeline").Parse(string(templateContent)))

	var output io.Writer

	if os.Getenv("STDOUT") != "" {
		output = os.Stdout
	} else {
		file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatal("Failed to open input file, reason is :", err)
		}

		defer file.Close()
		output = file
	}

	writer := bufio.NewWriter(output)

	for tweets := range out {
		if err := t.Execute(writer, tweets); err != nil {
			log.Println("Failed to execute template")
		}

		writer.Flush()
	}

	close(done)
}

func main() {
	client, err := setupClient(credentialsFilePath)

	if err != nil {
		log.Fatal("Failed to setup client, reason is", err)
	}

	c := getConfig(os.Args)

	in := make(chan *account, concurrency)
	out := make(chan []twitter.Tweet)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go fetchTimeline(client, in, out, &wg)
	}

	go dumpTimeline(c.OutputFile, out, done)

	loadAccounts(c.AccountsFile, in)

	close(in)

	wg.Wait()

	close(out)

	<-done

	log.Println("Done collecting tweets, enjoy :)")
}
