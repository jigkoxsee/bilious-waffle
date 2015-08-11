package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SlackPost struct {
	Token    string
	Team     string
	Channel  string
	Username string
	Text     string
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Home")
	fmt.Fprintf(w, "Home %s", r.URL.Path[1:])
}

func handlerThoth(w http.ResponseWriter, r *http.Request) {
	log.Println("Thoth")

	fbHost := os.Getenv("GO_HOST")
	slackToken := os.Getenv("GO_TOKEN")

	slack := SlackPost{
		Token:    r.FormValue("token"),
		Team:     r.FormValue("team_domain"),
		Channel:  r.FormValue("channel_name"),
		Username: r.FormValue("user_name"),
		Text:     r.FormValue("text"),
	}

	if slack.Token != slackToken {
		fmt.Fprintf(w, "Token not match")
	} else {
		fmt.Fprintf(w, "Thoth %s", slack)

		//--
		res, err := http.Get(fbHost + "/thoth/" + slack.Channel + ":" + slack.Username + ":" + slack.Text)
		if err != nil {
			log.Fatal(err)
		}

		bot, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Thoth req %s", bot)
	}

}

func main() {
	// Go Web
	port := ":" + os.Getenv("GO_PORT")

	// Slack
	slackToken := os.Getenv("GO_TOKEN")
	fmt.Println("Token:" + slackToken)

	// FB Chat
	// http://ec2-52-76-24-123.ap-southeast-1.compute.amazonaws.com:8000
	fbHost := os.Getenv("GO_HOST")
	fmt.Println("fbHost:" + fbHost)

	http.HandleFunc("/", handler)
	http.HandleFunc("/thoth", handlerThoth)
	http.ListenAndServe(port, nil)
}
