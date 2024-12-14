package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
