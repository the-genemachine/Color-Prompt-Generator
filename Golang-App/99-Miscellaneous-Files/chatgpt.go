package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func main() {
	// Replace YOUR_API_KEY with your actual OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": "Give me a tagline for a color prompt generator"}},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Failed to send the request: %v", err)
	}

	fmt.Println(response.String())
}
