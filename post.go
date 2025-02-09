package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Post(discordWebhook string, username string, content string) bool {
	client := new(http.Client)
	content, truncated := truncateContent(content)
	encodedContent, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}
	body_string := fmt.Sprintf(`{"username":"%s","content":%s}`, username, encodedContent)
	body := []byte(body_string)
	req, err := http.NewRequest("POST", discordWebhook, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		log.Println("Response status code is not 204")
	} else {
		log.Println("Successfully posted to discord")
	}
	return truncated
}
