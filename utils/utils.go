package utils

import (
	"log"
	"os"
)

var cwd string

func InitWD() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cwd = wd
}

func GetWd() string {
	if cwd == "" {
		InitWD()
	}

	return cwd
}
