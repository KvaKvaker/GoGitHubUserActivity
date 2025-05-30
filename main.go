package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var helpNote = `-------------------------------------------------------------
To get help print <file.exe> --help
To show github user’s recent activity print <file.exe> <username>
-------------------------------------------------------------
`

type ObjectJson struct {
	EventType string `json:"type"`
	Repo      struct {
		RepoId   int    `json:"id"`
		RepoName string `json:"name"`
		RepoUrl  string `json:"url"`
	}
	Payload struct {
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Branch  string `json:"master_branch"`
		Desc    string `json:"description"`
	}
	CreatedTime string `json:"created_at"`
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		fmt.Println(helpNote)
	} else {
		var ansObj []ObjectJson

		resp, err := http.Get("https://api.github.com/users/" + os.Args[1] + "/events/public")
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(data, &ansObj)

		if err != nil {
			log.Println(err)
			return
		}

		for _, el := range ansObj {
			if el.EventType == "CreateEvent" {
				if el.Payload.RefType == "repository" {
					fmt.Printf("Create repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
				} else if el.Payload.RefType == "branch" {
					fmt.Printf("Create branch '%v' in repo '%v' in %v\n", el.Payload.Ref, el.Repo.RepoName, el.CreatedTime)
				}
			} else if el.EventType == "PushEvent" {
				fmt.Printf("Push commit to repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else if el.EventType == "WatchEvent" {
				fmt.Printf("Watch for repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else if el.EventType == "PullRequestEvent" {
				fmt.Printf("Pull request in repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else if el.EventType == "PullRequestReviewEvent" {
				fmt.Printf("Pull request review in repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else if el.EventType == "IssuesEvent" {
				fmt.Printf("Opened a new issue in repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else if el.EventType == "IssueCommentEvent" {
				fmt.Printf("Add comment to issue in repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else {
				fmt.Println("\nUnknown type event :(\n)")
			}
		}
	}
}
