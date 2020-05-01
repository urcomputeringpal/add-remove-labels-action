package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v31/github"
	"github.com/sethvargo/go-githubactions"
)

func main() {

	// Validate input
	githubToken := githubactions.GetInput("GITHUB_TOKEN")
	if githubToken == "" {
		githubactions.Fatalf("missing 'githubToken'")
		return
	}

	labelString := githubactions.GetInput("labels")
	if labelString == "" {
		githubactions.Fatalf("missing 'labels'")
		return
	}
	labels := strings.Split(labelString, ",")

	action := githubactions.GetInput("action")
	if action == "" {
		githubactions.Fatalf("missing 'action'")
		return
	}

	// Setup GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Try to parse the event
	event, err := ioutil.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Printf("Can't read events: " + err.Error())
		return
	}

	var owner string
	var repo string
	var number int

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	switch eventName {
	case "issues":
		var issue github.IssuesEvent
		err = json.Unmarshal(event, &issue)
		if err != nil {
			log.Fatalf("Can't unmarshal json: %s", err)
		}

		owner = *issue.Repo.Owner.Login
		repo = *issue.Repo.Name
		number = *issue.Issue.Number
	case "pull_request":
		var pr github.PullRequestEvent
		err = json.Unmarshal(event, &pr)
		if err != nil {
			log.Fatalf("Can't unmarshal json: %s", err)
		}

		owner = *pr.Repo.Owner.Login
		repo = *pr.Repo.Name
		number = *pr.PullRequest.Number
	}

	log.Printf("%s/%s#%d %s %+v", owner, repo, number, action, labels)

	if action == "remove" {
		for _, label := range labels {
			_, err = client.Issues.RemoveLabelForIssue(ctx, owner, repo, number, label)
			if err != nil {
				log.Printf(err.Error())
				return
			}
		}
	} else {
		_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, number, labels)
		if err != nil {
			log.Printf(err.Error())
			return
		}
	}
}
