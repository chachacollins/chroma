package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const API_URL = "https://api.mistral.ai/v1/chat/completions"

var (
	star     = " You are one of the greatest programmers to ever live, you will receive code and your job would be to generate markdown documentation elaborating what the code does. You will return markdown and only markdown in the format specified by the starlight astro framework"
	markdown = " You are one of the greatest programmers to ever live, you will receive code and your job would be to generate markdown documentation elaborating what the code does. You will return markdown and only markdown.Make sure to keep your documentation brief but super clear "
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

func getApiKey() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: loading .env file: %s", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "API_KEY not found in environment variables")
	}
	return apiKey
}

func readFile(fileName string) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file")
	}
	return file
}

func generateMd(apiKey string, writeFilePath string, file []byte) {
	generateMdPayload := ChatRequest{
		Model: "codestral-latest",
		Messages: []Message{
			{
				Role:    "user",
				Content: string(file),
			},
			{
				Role:    "system",
				Content: markdown,
			},
		},
		Temperature: 0.5,
		MaxTokens:   100000,
	}

	payloadBytes, err := json.Marshal(generateMdPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: marshaling generateMd: %s", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: creating generateMd: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: making generateMd: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "Error: decoding response: %s", err)
	}

	if len(response.Choices) > 0 {
		writeFile(writeFilePath, response.Choices[0].Message.Content)
	}

}

func writeFile(filename string, data string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create file: %s", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not write to file: %s", err)
	}
	err = writer.Flush()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not flush buffer: %s", err)
	}

}
