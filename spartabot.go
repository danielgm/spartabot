package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	bieberLovePattern *regexp.Regexp
	slackToken        string
)

func main() {
	bieberLovePattern = regexp.MustCompile(`i (love|:heart:)[^.!?]*(justin)?(bieber|beiber)`)
	slackToken = os.Getenv("SLACK_TOKEN")
	log.Printf("Using Slack token: %s", slackToken)

	http.HandleFunc("/hook", hook)
	log.Println("Looking for Bieber love...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hook(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		msg := parseRequest(req)
		if msg != nil && msg["token"][0] == slackToken {
			lovesBieber := lovesJustinBieber(msg["text"][0])
			log.Printf("Love found! user=%s, channel=%s, bieber=%t, text=\"%s\"", msg["user_name"][0], msg["channel_name"][0], lovesBieber, msg["text"][0])
			if msg["user_name"][0] != "slackbot" && lovesBieber {
				fmt.Fprintf(res, "{\"text\": \"I love you, too, @%s.\"}", msg["user_name"][0])
			}
		}
	}
}

func parseRequest(req *http.Request) map[string][]string {
	b := new(bytes.Buffer)
	b.ReadFrom(req.Body)
	s := b.String()
	msg, err := url.ParseQuery(s)
	if err != nil {
		log.Printf("Bad webhook request. data=%s", s)
		return nil
	}
	return msg
}

func lovesJustinBieber(text string) bool {
	text = strings.ToLower(text)
	return bieberLovePattern.MatchString(text)
}
