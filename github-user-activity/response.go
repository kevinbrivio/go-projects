package main

import "time"

type GithubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type Actor struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
}

type Repo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	RepositoryID int64    `json:"repository_id"`
	PushID       int64    `json:"push_id"`
	Size         int      `json:"size"`
	DistinctSize int      `json:"distinct_size"`
	Ref          string   `json:"ref"`
	Head         string   `json:"head"`
	Before       string   `json:"before"`
	Action       string   `json:"action"`
	Commits      []Commit `json:"commits"`
}

type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
}

func (payload Payload) CommitCount() int {
	if payload.Size > 0 {
		return payload.Size
	}

	if payload.DistinctSize > 0 {
		return payload.DistinctSize
	}

	return len(payload.Commits)
}
