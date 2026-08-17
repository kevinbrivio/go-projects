package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const GITHUBURL = "https://api.github.com/users"
const ENDPOINT = "events"

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("github-activity ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		args := strings.Fields(line)

		if len(args) != 1 {
			fmt.Println("Use this format: [username]")
			fmt.Println("example: github-activity linus")
			continue
		}

		username := args[0]
		endpoint := fmt.Sprintf("%s/%s/%s", GITHUBURL, url.PathEscape(username), ENDPOINT)

		events, err := getUserEvents(endpoint)
		if err != nil {
			return
		}

		for _, event := range events {
			switch event.Type {
			case "PushEvent":
				commitCount := event.Payload.CommitCount()
				fmt.Printf("- Pushed %d commits to %s\n", commitCount, event.Repo.Name)

			case "IssuesEvent":
				if event.Payload.Action == "opened" {
					fmt.Printf("- Opened a new issue in %s\n", event.Repo.Name)
				}

			case "WatchEvent":
				if event.Payload.Action == "started" {
					fmt.Printf("- Starred %s\n", event.Repo.Name)
				}

			default:
				fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
			}
		}
	}
}

func getUserEvents(endpoint string) ([]GithubEvent, error) {
	body, err := getGithub(endpoint)
	if err != nil {
		return nil, err
	}

	var events []GithubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Couldn't parse JSON: ", err)
		return nil, err
	}

	return events, nil
}

func getGithub(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Couldn't prepare request:", err)
		return nil, err
	}

	if err := Load(); err != nil {
		fmt.Println("Couldn't load bearer token")
		return nil, err
	}
	githubToken := Token.GetValue()
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+githubToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error get: ", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Couldn't read reply: ", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github returned %s: %s", response.Status, string(body))
	}

	return body, nil
}
