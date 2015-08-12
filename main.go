package main

import (
	"fmt"
	"io/ioutil"
	//	"strconv"
	"encoding/json"
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

func fbBot(group string, msg string) ([]byte, error) {
	fbHost := os.Getenv("GO_HOST")
	res, err := http.Get(fbHost + "/" + group + "/" + msg)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return bot, err
}

func handlerThoth(w http.ResponseWriter, r *http.Request) {
	log.Println("Thoth")

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
		bot, err := fbBot("thoth", slack.Channel+":"+slack.Username+":"+slack.Text)
		if err != nil {
			log.Println(err)
		}
		//--
		fmt.Fprintf(w, "Thoth req %s", bot)
	}

}

func handlerLeafbox(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	//log.Println(dat)
	//log.Println(dat["action"].(string))

	// Case
	var msg string
	xEvent := r.Header.Get("X-GitHub-Event")
	switch xEvent {
	case "pull_request":
		msg = fmt.Sprintf("PR #%d %s by @%s", int(dat["number"].(float64)),
			dat["action"],
			dat["sender"].(map[string]interface{})["login"])
	case "pull_request_review_comment":
		msg = fmt.Sprintf("PR Comment #%d by @%s", int(dat["number"].(float64)),
			dat["sender"].(map[string]interface{})["login"])
	case "issues":
		msg = fmt.Sprintf("Issue #%d %s by @%s", int(dat["issue"].(map[string]interface{})["number"].(float64)),
			dat["action"],
			dat["sender"].(map[string]interface{})["login"])
	case "issue_comment":
		msg = fmt.Sprintf("Issue Comment #%d by @%s", int(dat["issue"].(map[string]interface{})["number"].(float64)),
			dat["sender"].(map[string]interface{})["login"])
	default:
		msg = ""
	}
	log.Println("MSG:" + msg)
	if msg != "" {
		bot, err := fbBot("leafbox", msg)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "Leafbox body %s", bot)
	} else {
		fmt.Fprintf(w, "Leafbox body %s", msg)
	}
	//fmt.Fprintf(w, "Leafbox json %s", t)
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
	http.HandleFunc("/leafbox", handlerLeafbox)
	http.ListenAndServe(port, nil)
}
