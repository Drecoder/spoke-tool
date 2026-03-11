package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/spoke-tool/internal/model"
)

func main() {
	// Create client
	client, err := model.NewClient(model.ClientConfig{
		OllamaHost: "http://localhost:11434",
		Timeout:    30 * time.Second,
		Models: map[model.ModelType]string{
			model.CodeLLamaEncoder: "codellama:7b",
			model.CodeLLamaDecoder: "codellama:7b",
			model.Gemma2BChat:      "gemma2:2b",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Test CodeLLama
	fmt.Println("Testing CodeLLama...")
	resp, err := client.Generate(context.Background(), model.SLMRequest{
		Model:     model.CodeLLamaEncoder,
		Prompt:    "Write a hello world function in Go",
		MaxTokens: 100,
	})
	if err != nil {
		log.Printf("CodeLLama error: %v", err)
	} else {
		fmt.Printf("CodeLLama response:\n%s\n", resp.Response)
	}

	// Test Gemma
	fmt.Println("\nTesting Gemma...")
	resp, err = client.Generate(context.Background(), model.SLMRequest{
		Model:     model.Gemma2BChat,
		Prompt:    "Explain what a hello world program does",
		MaxTokens: 100,
	})
	if err != nil {
		log.Printf("Gemma error: %v", err)
	} else {
		fmt.Printf("Gemma response:\n%s\n", resp.Response)
	}
}
