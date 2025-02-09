package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func PostFile(discordWebhook string, username string, filePath string, message string) {
	content, _ := truncateContent(message)
	client := new(http.Client)
	b, w := createMultipartFormDataWithPayload("Document", filePath, username, fmt.Sprintf("%s attached", content))

	req, err := http.NewRequest("POST", discordWebhook, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		log.Printf("Uploaded %s to discord. Response status code is not 204 (status code %d)", message, resp.StatusCode)
	} else {
		log.Printf("Successfully uploaded %s to discord", message)
	}
}

func createMultipartFormDataWithPayload(fieldName string, fileName string, username string, content string) (bytes.Buffer, *multipart.Writer) {
	var buffer bytes.Buffer
	var err error
	writer := multipart.NewWriter(&buffer)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = writer.CreateFormFile(fieldName, file.Name()); err != nil {
		log.Fatalf("Error creating writer: %v", err)
	}
	defer file.Close()
	if _, err = io.Copy(fw, file); err != nil {
		log.Fatalf("Error with io.Copy: %v", err)
	}

	encodedContent, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}
	fieldWriter, err := writer.CreateFormField("payload_json")
	payload := fmt.Sprintf(`{"username":"%s","content":%s}`, username, encodedContent)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fieldWriter.Write([]byte(payload))
	if err != nil {
		log.Fatal(err)
	}

	writer.Close()
	return buffer, writer
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}
