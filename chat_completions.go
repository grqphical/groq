package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Conversation is a struct that allows you to construct chat completion requests
type Conversation struct {
	messages     []Message
	systemPrompt string
}

// Creates a new conversation
func NewConversation(systemPrompt string) *Conversation {
	var messages []Message = make([]Message, 0)

	messages = append(messages, Message{
		Role:    MessageRoleSystem,
		Content: systemPrompt,
	})

	return &Conversation{
		messages,
		systemPrompt,
	}
}

// adds message(s) to the conversations
func (c *Conversation) AddMessages(messages ...Message) {
	c.messages = append(c.messages, messages...)
}

// clears the history of the conversation
func (c *Conversation) ClearHistory() {
	c.messages = []Message{
		{
			Role:    MessageRoleSystem,
			Content: c.systemPrompt,
		},
	}
}

// sends the conversation to the AI returning the API's result and adding the message to the conversation's history
func (c *Conversation) Complete(g *GroqClient, model string, config *ChatCompletionConfig) (ChatCompletionResponse, error) {
	var requestBody ChatCompletionRequest

	requestBody.Model = model
	requestBody.Messages = c.messages

	if config != nil {
		requestBody.Temperature = config.Temperature
		requestBody.TopP = config.TopP
		requestBody.Stream = config.Stream
		requestBody.Stop = config.Stop
		requestBody.MaxTokens = config.MaxTokens
		requestBody.PresencePenalty = config.PresencePenalty
		requestBody.FrequencyPenalty = config.FrequencyPenalty
		requestBody.User = config.User
	} else {
		requestBody.ResponseFormat.Type = "text"
	}

	marshaledJson, err := json.Marshal(requestBody)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req, err := createGroqRequest("chat/completions", g.apiKey, "POST", bytes.NewBuffer(marshaledJson))
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	if resp.StatusCode != 200 {
		return ChatCompletionResponse{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var chatResponse ChatCompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	c.AddMessages(Message{
		Role:    MessageRoleAssistant,
		Content: chatResponse.Choices[0].Message.Content,
	})

	return chatResponse, nil
}

type ChatCompletionRequest struct {
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`

	Messages []Message `json:"messages"`

	Model           string  `json:"model"`
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	ResponseFormat struct {
		Type string `json:"type,omitempty"`
	} `json:"response_format,omitempty"`

	Seed        int      `json:"seed,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	User        string   `json:"user,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`
}

type ChatCompletionConfig struct {
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float64

	// Maximum amount of tokens that can be generated in the completion
	MaxTokens int

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	PresencePenalty float64

	ResponseFormat struct {
		Type string
	}

	// Random seed for the model
	Seed int

	// Up to 4 sequences where the API will stop generating tokens
	Stop []string

	// Whether or not the API should stream responses. Currently UNSUPPORTED
	Stream bool
	// The sampling temperature, between 0 and 1.
	Temperature float64

	// Unique identifier for your end-user
	User string

	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// DO NOT altering this if you have altered temperature and vice versa.
	TopP float64
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		Logprobs     any     `json:"logprobs"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		QueueTime        float64 `json:"queue_time"`
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"x_groq"`
}
