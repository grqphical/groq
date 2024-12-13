package groq_test

import (
	"os"
	"testing"

	"github.com/grqphical/groq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroqClient(t *testing.T) {
	invalidApiKey := "notAnApiKey"

	_, err := groq.NewGroqClient(invalidApiKey)
	assert.Error(t, err, "api key should be invalid")
}

func TestGetModels(t *testing.T) {
	godotenv.Load(".env")

	apiKey := os.Getenv("GROQ_TOKEN")
	if apiKey == "" {
		assert.Fail(t, "API key not set. Make sure you have an environment vairable (or .env file) with GROQ_TOKEN set to your API key")
		return
	}

	client, err := groq.NewGroqClient(apiKey)
	assert.NoError(t, err, "the API key should be valid. double check you have set the environment variable correctly")

	models, err := client.GetModels()
	assert.NoError(t, err, "the request should succeed")
	assert.Greater(t, len(models), 0, "there should be models returned")
}

func TestGetModel(t *testing.T) {
	godotenv.Load(".env")

	apiKey := os.Getenv("GROQ_TOKEN")
	if apiKey == "" {
		assert.Fail(t, "API key not set. Make sure you have an environment vairable (or .env file) with GROQ_TOKEN set to your API key")
		return
	}

	client, err := groq.NewGroqClient(apiKey)
	assert.NoError(t, err, "the API key should be valid. double check you have set the environment variable correctly")

	model, err := client.GetModel("llama3-8b-8192")
	assert.NoError(t, err, "the request should succeed")
	assert.Equal(t, "llama3-8b-8192", model.Id)
	assert.Equal(t, "Meta", model.OwnedBy)
	assert.Equal(t, 1693721698, model.Created)
}
