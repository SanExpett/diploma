package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/SanExpett/diploma/internal/domain"
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

	err = insertUsers()
	if err != nil {
		log.Printf("failed to insert users data: %v \n", err)

		return err
	}

	insertComments()

	return nil
}

func insertInternal() error {
	jsonData, err := os.ReadFile("cmd/scripts/data/internal.json")
	if err != nil {
		return err
	}

	var filmsToAdd []*domain.FilmToAdd

	err = json.Unmarshal(jsonData, &filmsToAdd)
	if err != nil {
		return err
	}

	for _, filmToAdd := range filmsToAdd {
		filmB, err := json.Marshal(&filmToAdd)
		if err != nil {
			return err
		}

		resp, err := http.Post("http://127.0.0.1:8081/api/films/add", "application/json", bytes.NewBuffer(filmB))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertUsers() error {
	jsonData, err := os.ReadFile("cmd/scripts/data/users.json")
	if err != nil {
		log.Printf("failed to read users data: %v \n", err)
		return err
	}

	var users []json.RawMessage
	if err := json.Unmarshal(jsonData, &users); err != nil {
		log.Printf("failed to parse users data: %v \n", err)
		return err
	}

	for _, userData := range users {
		resp, err := http.Post("http://127.0.0.1:8081/api/auth/signup", "application/json", bytes.NewBuffer(userData))
		if err != nil {
			log.Printf("failed to send signup request: %v \n", err)
			return err
		}
		defer resp.Body.Close()

		log.Printf("resp cookies: %v \n", resp.Cookies())

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read response: %v \n", err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("signup failed with status %d: %s\n", resp.StatusCode, string(body))
			return fmt.Errorf("signup failed: %s", string(body))
		}
	}

	return nil
}

func insertComments() {
	jsonData, err := os.ReadFile("cmd/scripts/data/users.json")
	if err != nil {
		log.Printf("failed to read users data: %v \n", err)

		return
	}

	var users []struct {
		Email    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(jsonData, &users); err != nil {
		log.Printf("failed to parse users data: %v \n", err)

		return
	}

	comments := []string{
		"Great film!",
		"Amazing story",
		"Interesting plot",
		"Would recommend to watch",
		"Nice movie experience",
		"Excellent performance by the actors",
		"Masterpiece of cinematography",
		"Captivating from start to finish",
		"A must-see film",
		"Unforgettable viewing experience",
	}

	for _, user := range users {
		// Login to get access cookie
		loginData := map[string]string{
			"login":    user.Email,
			"password": user.Password,
		}
		loginJSON, err := json.Marshal(loginData)
		if err != nil {
			log.Printf("failed to marshal login data: %v\n", err)
			continue
		}

		loginResp, err := http.Post("http://127.0.0.1:8081/api/auth/login", "application/json", bytes.NewBuffer(loginJSON))
		if err != nil {
			log.Printf("failed to login: %v\n", err)
			continue
		}
		defer loginResp.Body.Close()

		log.Printf("loginResp cookies: %v\n", loginResp.Cookies())

		if loginResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(loginResp.Body)
			log.Printf("login failed with status %d: %s\n", loginResp.StatusCode, string(body))

			continue
		}

		// Get access cookie
		var accessCookie, uuidCookie *http.Cookie
		for _, cookie := range loginResp.Cookies() {
			if cookie.Name == "access" {
				accessCookie = cookie
			}

			if cookie.Name == "user_uuid" {
				uuidCookie = cookie
			}

			if accessCookie != nil && uuidCookie != nil {
				break
			}
		}

		if accessCookie == nil {
			log.Println("access cookie not found")

			continue
		}

		if uuidCookie == nil {
			log.Println("uuid cookie not found")

			continue
		}

		// Get films
		client := &http.Client{}
		filmsReq, err := http.NewRequest("GET", "http://127.0.0.1:8081/api/films/all", nil)
		if err != nil {
			log.Printf("failed to create films request: %v\n", err)

			continue
		}
		filmsReq.AddCookie(accessCookie)

		filmsResp, err := client.Do(filmsReq)
		if err != nil {
			log.Printf("failed to get films: %v\n", err)

			continue
		}
		defer filmsResp.Body.Close()

		resp := domain.FilmsPreviewsResponse{}
		if err := json.NewDecoder(filmsResp.Body).Decode(&resp); err != nil {
			log.Printf("failed to decode films: %v\n", err)

			continue
		}

		for _, film := range resp.Films {
			commentData := &domain.CommentToAdd{
				AuthorUuid: uuidCookie.Value,
				FilmUuid:   film.Uuid,
				Text:       comments[rand.Intn(len(comments))],
				Score:      uint32(rand.Intn(5) + 1),
			}

			commentJSON, err := json.Marshal(commentData)
			if err != nil {
				log.Printf("failed to marshal comment: %v\n", err)

				continue
			}

			commentReq, err := http.NewRequest("POST", "http://127.0.0.1:8081/api/films/comments/add", bytes.NewBuffer(commentJSON))
			if err != nil {
				log.Printf("failed to create comment request: %v\n", err)

				continue
			}
			commentReq.AddCookie(accessCookie)
			commentReq.Header.Set("Content-Type", "application/json")

			commentResp, err := client.Do(commentReq)
			if err != nil {
				log.Printf("failed to post comment: %v\n", err)

				continue
			}
			defer commentResp.Body.Close()

			if commentResp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(commentResp.Body)
				log.Printf("comment failed with status %d: %s\n", commentResp.StatusCode, string(body))

				continue
			}
		}
	}
}
