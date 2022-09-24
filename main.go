package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// ACCESS_TOKEN use to access Telegram API
	ACCESS_TOKEN = "5756171235:AAH3x6HTAMg5iHJ0YPCRTX8jjuAk3RY25fI"
	// TELEGRAM_API_URL is a base URL for Telegram API
	TELEGRAM_API_URL = "https://api.telegram.org"
)

// Message data structure for telegram message
type Message struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"message"`
}

// webhook is a handler for Webhook server
func webhook(w http.ResponseWriter, r *http.Request) {
	// return all with status code 200
	w.WriteHeader(http.StatusOK)

	// read body in the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
	}

	// initiate Message data structure to message variable
	// unmarshal []byte data into message
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("failed to unmarshal body: %v", err)
		return
	}

	// send message to end-user
	err = sendMessage(message.Message.Chat.ID, "Automatically Reply 🙌🏻")
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}

	return
}

// sendMessage sends a message to end-user
func sendMessage(chatID int, message string) error {
	// setup http request
	url := fmt.Sprintf("%s/bot%s/sendMessage?chat_id=%d&text=%s", TELEGRAM_API_URL, ACCESS_TOKEN, chatID, message)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}

	// send http request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	// print response
	log.Printf("message sent successfully?\n%#v", res)

	return nil
}

func main() {
	// create the handler
	handler := http.NewServeMux()
	handler.HandleFunc("/", webhook)

	// configure http server
	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("localhost:%d", 3000),
	}

	// start http server
	log.Printf("http server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
