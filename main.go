package main

import (
	"flag"
	"log"
)

func main() {
	var seedApiToken, slackToken string
	flag.StringVar(&seedApiToken, "seed-token", "", "Seed api token")
	flag.StringVar(&slackToken, "slack-token", "", "Slack integration token")
	flag.Parse()

	if seedApiToken == "" {
		log.Fatalf("seed-token is required")
	}

	if slackToken == "" {
		log.Fatalf("slack-token is required")
	}

	sb := New(slackToken, seedApiToken)
	sb.Run()
}
