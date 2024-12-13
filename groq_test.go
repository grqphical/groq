package groq_test

import (
	"testing"

	"github.com/grqphical/groq"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroqClient(t *testing.T) {
	invalidApiKey := "notAnApiKey"

	_, err := groq.NewGroqClient(invalidApiKey)
	assert.Error(t, err, "api key should be invalid")
}
