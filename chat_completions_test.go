package groq_test

import (
	"os"
	"testing"

	"github.com/grqphical/groq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestChatCompletion(t *testing.T) {
	godotenv.Load(".env")

	apiKey := os.Getenv("GROQ_TOKEN")
	if apiKey == "" {
		assert.Fail(t, "API key not set. Make sure you have an environment vairable (or .env file) with GROQ_TOKEN set to your API key")
		return
	}

	client, err := groq.NewGroqClient(apiKey)
	assert.NoError(t, err, "the API key should be valid. double check you have set the environment variable correctly")

	conversation := groq.NewConversation("")

	conversation.AddMessages(groq.Message{
		Role:    groq.MessageRoleUser,
		Content: "say the word banana in lowercase with no punctuation",
	})

	response, err := conversation.Complete(client, "llama3-8b-8192", nil)
	assert.NoError(t, err)
	assert.Equal(t, "banana", response.Choices[0].Message.Content, "should've said banana")
}
