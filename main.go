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

		resp, err := http.Get("https://api.github.com/users/KvaKvaker/events/public")
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close() // закрываем тело ответа после работы с ним

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
			//fmt.Printf("id - %v\t el - %v\n", idx, el)
			if el.EventType == "CreateEvent" {
				if el.Payload.RefType == "repository" {
					fmt.Printf("Create repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
				} else if el.Payload.RefType == "branch" {
					fmt.Printf("Create branch '%v' in repo '%v' in %v\n", el.Payload.Ref, el.Repo.RepoName, el.CreatedTime)
				}
			} else if el.EventType == "PushEvent" {
				fmt.Printf("Push commit to repo '%v' in %v\n", el.Repo.RepoName, el.CreatedTime)
			} else {
				fmt.Println("\nUnknown type event :(\n)")
			}
		}

		//fmt.Printf("%s", ansObj) // печатаем ответ как строку

	}
}
