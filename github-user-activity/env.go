package main

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

const (
	Token EnvKey = "GITHUB_TOKEN"
)

func Load() error {
	return godotenv.Load(".env")
}
