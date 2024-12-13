// Package groq is an unofficial API wrapper for GroqCloud https://groq.com/
//
// groq requires Go 1.14 or newer
package groq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// creates an http.Request with the API key added to it and the URL set
func createGroqRequest(endpoint string, apiKey string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.groq.com/openai/v1/%s", endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", apiKey))

	return req, nil
}

// Creates a new Groq client. Returns an error if the API key given is invalid
func NewGroqClient(apiKey string) (*GroqClient, error) {
	// test the API key
	req, err := createGroqRequest("/models", apiKey)
	if err != nil {
		return nil, err
	}

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

// returns all models available on GroqCloud
func (g *GroqClient) GetModels() ([]Model, error) {
	req, err := createGroqRequest("/models", g.apiKey)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	var response modelsResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
