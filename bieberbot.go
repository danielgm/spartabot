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
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(req.Body)
		msg, err := url.ParseQuery(buffer.String())
		if err != nil {
			panic(err)
		}

		fmt.Println("user ID: " + msg["user_id"][0])
		fmt.Println("user name: " + msg["user_name"][0])
		if lovesJustinBieber(msg["text"][0]) {
			log.Printf(
				"#%s user:%s (%s), \"%s\" (\"%s\")\n",
				msg["channel_name"][0], msg["user_id"][0], msg["user_name"][0], msg["text"][0], msg["trigger_word"][0])
			fmt.Fprintf(res, "{\"text\": \"Oh, I love you, too, @%s.\"}", msg["user_name"][0])
		}
	}
}

func lovesJustinBieber(text string) bool {
	text = strings.ToLower(text)
	return bieberLovePattern.MatchString(text)
}
