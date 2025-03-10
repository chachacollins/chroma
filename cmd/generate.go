package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chachacollins/chroma/cfg"
	"github.com/chachacollins/chroma/utils"
)

const API_URL = "https://api.mistral.ai/v1/chat/completions"

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

func readFile(fileName string) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file")
	}
	return file
}

func generateMd(writeFilePath string, file []byte) {
	utils.PrintLog("Generating markdown...")
	generateMdPayload := ChatRequest{
		Model: "codestral-latest",
		Messages: []Message{
			{
				Role:    "user",
				Content: string(file),
			},
			{
				Role:    "system",
				Content: cfg.Load().MarkdownPrompt,
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
		fmt.Fprintf(os.Stderr, "Error: creating generateMd: %s\n", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Load().Apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: making generateMd: %s\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "unexpected status code: %d\n", resp.StatusCode)
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "Error: decoding response: %s\n", err)
	}

	if len(response.Choices) > 0 {
		utils.WriteFile(writeFilePath, response.Choices[0].Message.Content)
	}
	utils.PrintSuccess("Generated markdown")
}

func inlineComm(writeFilePath string, file []byte) {
	utils.PrintLog("Generating inline documentation...")
	inlineCommPayload := ChatRequest{
		Model: "codestral-latest",
		Messages: []Message{
			{
				Role:    "user",
				Content: string(file),
			},
			{
				Role:    "system",
				Content: cfg.Load().InlinePrompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   100000,
	}

	payloadBytes, err := json.Marshal(inlineCommPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: marshaling inlineComm: %s\n", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: creating inlineComm: %s\n", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Load().Apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: making inlineComm: %s\n", err)
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
		fmt.Fprintf(os.Stderr, "Error: decoding response: %s\n", err)
	}

	if len(response.Choices) > 0 {
		content := utils.ExtractCodeBlocks(response.Choices[0].Message.Content)
		utils.WriteFile(writeFilePath, content[0])
	}
	utils.PrintSuccess("Generated inline documentation")

}

func starLight(writeFilePath string, file []byte) {
	utils.PrintLog("Generating starlight documentation...")
	starLightPayload := ChatRequest{
		Model: "codestral-latest",
		Messages: []Message{
			{
				Role:    "user",
				Content: string(file),
			},
			{
				Role:    "system",
				Content: cfg.Load().StarPrompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   10000,
	}

	payloadBytes, err := json.Marshal(starLightPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: marshaling starLight: %s\n", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: creating starLight: %s\n", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Load().Apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: making starLight: %s\n", err)
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
		fmt.Fprintf(os.Stderr, "Error: decoding response: %s\n", err)
	}
	fileNew := filepath.Base(writeFilePath)
	filePath := "./docs/src/content/docs/reference/" + fileNew + ".md"

	if len(response.Choices) > 0 {
		utils.WriteFile(filePath, response.Choices[0].Message.Content)
	}
	utils.PrintSuccess("Generated starlight docs")
}

func getDocs() {
	utils.PrintLog("Fetching docs...")
	_, errr := os.ReadDir("./docs")
	if errr != nil {
		cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/chachacollins/chromatemplate.git", "docs")
		_, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not run command: %s\n", err)
		}
		cmd = exec.Command("rm", "-rf", ".git")
		cmd.Dir = "./docs"
		_, err = cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not run command: %s\n", err)
		}
		utils.PrintSuccess("Finished getting docs")
		installDocs()
	}
}

func installDocs() {
	utils.PrintLog("installing docs...")
	cmd := exec.Command("npm", "i")
	cmd.Dir = "./docs"
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not run command: %s", err)
	}
	utils.PrintSuccess("Installed docs successfully")
}

func serveDocs() {
	_, err := os.ReadDir("docs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not find docs directory: %s\n", err)
		os.Exit(1)
	}
	buildDocs()
	runPreviewDocs()
}
func buildDocs() {
	utils.PrintLog("Building the docs...")
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "docs"
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not run command: %s\n", err)
	}
	utils.PrintSuccess("Finished building docs")

}

func runPreviewDocs() {
	cmd := exec.Command("npm", "run", "preview")
	cmd.Dir = "docs"
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not start command: %s\n", err)
		return
	}

	utils.PrintLog("running the server at localhost:4321")

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: command exited with error: %s\n", err)
	}
}
