package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var grokKey = `<Your API Key Here>`
var baseURL = "https://api.x.ai/v1"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func TextChat(tosend string) {
	tosend += ". Just Provide The code"
	// Create the request payload
	payload := ChatCompletionRequest{
		Model: "grok-2-1212",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: tosend,
			},
		},
	}

	// Marshal the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", baseURL+"/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+grokKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned non-200 status: %d\nResponse: %s\n", resp.StatusCode, string(body))
		return
	}

	var chatResp ChatCompletionResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	// Extract and print the response
	response := chatResp.Choices[0].Message.Content

	if !strings.Contains(response, "```") {
		fmt.Println(response)
		return
	}

	flag := false
	parts := strings.Split(response, "\n")
	for _, part := range parts {
		if strings.HasPrefix(part, "```") && !flag {
			flag = true
			continue
		} else if strings.HasPrefix(part, "```") && flag {
			flag = false
			continue
		}

		if flag {
			fmt.Println(part)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + ` "Can you please generate an ansible playbook for tagging a virtual machine on vmware"`)
		return
	}

	tosend := os.Args[1]

	TextChat(tosend)
}

