package groq

import (
	"errors"
	"fmt"
	"net/http"
)

// GroqClient is the main client that interacts with the GroqCloud API
type GroqClient struct {
	apiKey string
}

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

	return &GroqClient{
		apiKey,
	}, nil
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
	Active        bool `json:"active"`
	ContextWindow int  `json:"context_window"`
}
