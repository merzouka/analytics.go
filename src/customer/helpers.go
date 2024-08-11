package main

import (
	"fmt"
	"log"
	"os"
)

func setLogger() *os.File {
	path := os.Getenv("LOGS_PATH")
	if path == "" {
		path = "./logs"
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		log.Fatal("failed to open logs file")
	}

	prefix := os.Getenv("SERVICE_NAME")
	if prefix != "" {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	log.SetPrefix(prefix)
	log.SetOutput(f)
	return f
}

func GetUrl(endpoint string) string {
    return fmt.Sprintf("%s/%s", os.Getenv("TRANSACTION_SERVICE"), endpoint)
}
