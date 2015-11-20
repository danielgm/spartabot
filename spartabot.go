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
	patternResponseMap map[*regexp.Regexp]string
	slackToken         string
)

func main() {
	patternResponseMap[regexp.MustCompile(`spartans`)] = "Awoo! Awoo! Awoo!"
	patternResponseMap[regexp.MustCompile(`What is your profession`)] = "Awoo! Awoo! Awoo!"
	patternResponseMap[regexp.MustCompile(`Respect and honor`)] = "Respect and honor"
	patternResponseMap[regexp.MustCompile(`Respect and honour`)] = "Respect and honour"
	patternResponseMap[regexp.MustCompile(`This is madness`)] = "Madness? This is Sparta!"
	patternResponseMap[regexp.MustCompile(`Give them nothing`)] = "But take from them everything!"
	patternResponseMap[regexp.MustCompile(`Our arrows will blot out the sun`)] = "Then we will fight in the shade!"
	patternResponseMap[regexp.MustCompile(`There is much our cultures could share`)] = "Haven't you noticed? We've been sharing our culture with you all morning."

	slackToken = os.Getenv("SLACK_TOKEN")
	log.Printf("Using Slack token: %s", slackToken)

	http.HandleFunc("/hook", hook)
	log.Println("Listening for Spartans...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hook(res http.ResponseWriter, req *http.Request) {
	if isValidRequest(req) {
		msg := parseRequest(req)
		if isValidMessage(msg) {
			requestText := getMessageText(msg)
			responseText := getResponseText(requestText)
			if len(responseText) > 0 {
				log.Printf("Matched! user=%s, channel=%s, text=\"%s\"", msg["user_name"][0], msg["channel_name"][0], msg["text"][0])
				fmt.Fprintf(res, "{\"text\": \"%s\"}", responseText)
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

func isValidRequest(req *http.Request) bool {
	return req.Method == "POST"
}

func isValidMessage(msg map[string][]string) bool {
	return msg != nil && msg["token"][0] == slackToken
}

func getMessageText(msg map[string][]string) string {
	return msg["text"][0]
}

func getResponseText(msg string) string {
	msg = strings.ToLower(msg)
	for pattern, response := range patternResponseMap {
		if pattern.MatchString(msg) {
			return response
		}
	}
	return ""
}
