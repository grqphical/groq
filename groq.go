// Package groq is an unofficial API wrapper for GroqCloud https://groq.com/
//
// groq requires Go 1.14 or newer
package groq

import (
	"errors"
	"fmt"
	"net/http"
)

// Creates a new Groq client. Returns an error if the API key given is invalid
func NewGroqClient(apiKey string) (*GroqClient, error) {
	// test the API key
	req, err := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// if the API key is invalid return an error
	if resp.StatusCode == 401 {
		return nil, errors.New("invalid API key")
	}

	// return the client
	return &GroqClient{
		apiKey,
	}, nil
}
