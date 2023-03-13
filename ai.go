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

	file, err := os.OpenFile("interaction.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)

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

		payload := request{
			Model:       aiModel,
			Messages:    []message{{Role: "user", Content: escapedInput}},
			Temperature: temperature,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling payload: %s\n", err.Error())
			logger.Printf("Error marshalling payload: %s\n", err.Error())
			continue
		}

		req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Printf("Error creating request: %s\n", err.Error())
			logger.Printf("Error creating request: %s\n", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request: %s\n", err.Error())
			logger.Printf("Error sending request: %s\n", err.Error())
			continue
		}
		defer resp.Body.Close()

		var response response
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Printf("Error decoding response: %s\n", err.Error())
			logger.Printf("Error decoding response: %s\n", err.Error())
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error: unexpected status code %d", resp.StatusCode)
			logger.Printf("Error: unexpected status code %d", resp.StatusCode)
			continue
		}

		if len(response.Choices) == 0 {
			log.Println("Error: empty response")
			logger.Println("Error: empty response")
			continue
		}

		var message message
		if response.Choices[0].Message.Content != "" {
			message = response.Choices[0].Message
			logger.Printf("AI: %s\n", message.Content)
			log.Printf("AI: %s\n", message.Content)
		}
	}
}
