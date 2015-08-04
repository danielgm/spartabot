package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SlackMessage struct {
	token        string
	team_id      string
	team_domain  string
	channel_id   string
	channel_name string
	timestamp    float64
	user_id      string
	user_name    string
	text         string
	trigger_word string
}

var bieberLovePattern *regexp.Regexp

func main() {
	bieberLovePattern = regexp.MustCompile(`i love[^.!?]*(justin)?bieber`)

	http.HandleFunc("/hook", hook)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hook(res http.ResponseWriter, req *http.Request) {
	fmt.Println("hooked")
	if req.Method == "POST" {
		msg := parseSlackMessage(req.Body)
		fmt.Println("user ID: " + msg.user_id)
		fmt.Println("user name: " + msg.user_name)
		if lovesJustinBieber(msg) {
			log.Printf(
				"#%s user:%s (%s), \"%s\" (\"%s\")\n",
				msg.channel_name, msg.user_id, msg.user_name, msg.text, msg.trigger_word)
			fmt.Fprintf(res, "{\"text\": \"Oh, I love you, too, %s.\"}", msg.user_name)
		}
	}
}

func parseSlackMessage(body io.Reader) SlackMessage {
	var lines []string

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	m := make(map[string]string)
	for _, line := range lines {
		fmt.Println("> " + line)
		tokens := strings.Split(line, "=")
		if len(tokens) > 1 {
			key, value := tokens[0], tokens[1]
			m[key] = value
		}
	}

	timestamp, err := strconv.ParseFloat(m["timestamp"], 64)
	if err != nil {
		timestamp = 0
	}

	return SlackMessage{
		m["token"],
		m["team_id"],
		m["team_domain"],
		m["channel_id"],
		m["channel_name"],
		timestamp,
		m["user_id"],
		m["user_name"],
		m["text"],
		m["trigger_word"],
	}
}

func lovesJustinBieber(msg SlackMessage) bool {
	text := strings.ToLower(msg.text)
	return bieberLovePattern.MatchString(text)
}
