package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiEndpoint = "https://api.openai.com/v1/chat/completions"
const temperature = 0.5
const aiModel = "gpt-3.5-turbo-0301"

type request struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	logger := createLogger("interaction.log")

	client := &http.Client{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		escapedInput := strconv.Quote(input)
		logger.Printf(input)
		if strings.ToLower(input) == "exit" {
			break
		}

		response, err := getAIResponse(apiKey, client, escapedInput)
		if err != nil {
			logger.Printf("Error getting AI response: %s\n", err.Error())
			continue
		}

		if len(response.Choices) == 0 {
			logger.Println("Error: empty response")
			continue
		}

		message := response.Choices[0].Message
		if message.Content != "" {
			logAndPrintResponse(logger, message)
		}
	}
}

func createLogger(logFile string) *log.Logger {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "", log.LstdFlags)
}

func getAIResponse(apiKey string, client *http.Client, escapedInput string) (*response, error) {
	payload := request{
		Model:       aiModel,
		Messages:    []message{{Role: "user", Content: escapedInput}},
		Temperature: temperature,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %s", err.Error())
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseObj response
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	return &responseObj, nil
}

func logAndPrintResponse(logger *log.Logger, message message) {
	logger.Printf("AI: %s\n", message.Content)
	log.Printf("AI: %s\n", message.Content)
}
