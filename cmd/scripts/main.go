package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := insert(); err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}

func insert() error {
	err := insertInternal()
	if err != nil {
		log.Printf("failed to insert internal data: %v \n", err)
		return err
	}

	return nil
}

func insertInternal() error {
	jsonData, err := os.ReadFile("data/internal.json")
	if err != nil {
		return err
	}

	resp, err := http.Post("http://127.0.0.1:8081/api/films/add", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
