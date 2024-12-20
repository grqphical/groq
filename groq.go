// Package groq is an unofficial API wrapper for GroqCloud https://groq.com/
//
// groq requires Go 1.14 or newer
package groq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	MessageRoleUser      = "user"
	MessageRoleSystem    = "system"
	MessageRoleAssistant = "assistant"
)

// GroqClient is the main client that interacts with the GroqCloud API
type GroqClient struct {
	apiKey string
}

// a struct that represents the exact response returned by Groq's API
type modelsResponse struct {
	Data []Model `json:"data"`
}

// Model represents the metadata for an LLM hosted on Groqcloud
type Model struct {
	// the model's ID
	Id     string `json:"id"`
	Object string `json:"object"`

	// Unix timestamp when the model was created
	Created int `json:"created"`

	// Who owns this model
	OwnedBy string `json:"owned_by"`

	// Is the model currently active?
	Active bool `json:"active"`

	// How many context window tokens the model supports
	ContextWindow int `json:"context_window"`
}

// creates an http.Request with the API key added to it and the URL set
func createGroqRequest(endpoint string, apiKey string, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://api.groq.com/openai/v1/%s", endpoint), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return req, nil
}

// Creates a new Groq client. Returns an error if the API key given is invalid
func NewGroqClient(apiKey string) (*GroqClient, error) {
	// test the API key
	req, err := createGroqRequest("/models", apiKey, "GET", nil)
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
	req, err := createGroqRequest("/models", g.apiKey, "GET", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var response modelsResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// returns the information for a specific model
func (g *GroqClient) GetModel(modelId string) (Model, error) {
	req, err := createGroqRequest(fmt.Sprintf("/models/%s", modelId), g.apiKey, "GET", nil)
	if err != nil {
		return Model{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Model{}, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			return Model{}, errors.New("invalid model id")
		}
		return Model{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var model Model

	err = json.NewDecoder(resp.Body).Decode(&model)
	if err != nil {
		return Model{}, err
	}

	return model, nil
}
