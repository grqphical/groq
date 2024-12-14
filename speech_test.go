package groq_test

import (
	"os"
	"testing"

	"github.com/grqphical/groq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestTranscription(t *testing.T) {
	godotenv.Load(".env")

	apiKey := os.Getenv("GROQ_TOKEN")
	if apiKey == "" {
		assert.Fail(t, "API key not set. Make sure you have an environment vairable (or .env file) with GROQ_TOKEN set to your API key")
		return
	}

	client, err := groq.NewGroqClient(apiKey)
	assert.NoError(t, err, "the API key should be valid. double check you have set the environment variable correctly")

	assert.NoError(t, err)

	transcription, err := client.TranscribeAudio("testdata/testAudio.mp3", "whisper-large-v3", nil)
	assert.NoError(t, err)
	assert.Equal(t, " The quick brown fox jumped over the lazy dog.", transcription.Text)

	// test text responses
	transcription, err = client.TranscribeAudio("testdata/testAudio.mp3", "whisper-large-v3", &groq.TranscriptionConfig{
		ResponseFormat: "text",
	})
	assert.NoError(t, err)
	assert.Equal(t, " The quick brown fox jumped over the lazy dog.", transcription.Text)
}
