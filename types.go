package groq

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

type transcriptionSegment struct {
	ID                string  `json:"id"`
	Seek              float64 `json:"seek"`
	Start             float64 `json:"start"`
	End               float64 `json:"end"`
	Text              string  `json:"text"`
	Tokens            []int   `json:"tokens"`
	Temperature       int     `json:"temperature"`
	AvgLogProb        float64 `json:"avg_logprob"`
	CompressionRation float64 `json:"compression_ratio"`
	NoSpeechProb      float64 `json:"no_speech_prob"`
}

// TranscriptionConfig houses configuration options for transcription requests
type TranscriptionConfig struct {
	Language       string  `json:"language"`
	Prompt         string  `json:"prompt"`
	ResponseFormat string  `json:"response_format"`
	Temperature    float64 `json:"temperature"`
}

// this internal struct is used to create the request body for transcriptions
type transcriptionRequest struct {
	File           string  `json:"file"`
	Language       string  `json:"language"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt"`
	ResponseFormat string  `json:"response_format"`
	Temperature    float64 `json:"temperature"`
}

// Transcription represents an audio transcription result from one of Groq's models
type Transcription struct {
	Task     string                 `json:"task"`
	Language string                 `json:"language"`
	Duration float64                `json:"duration"`
	Text     string                 `json:"text"`
	Segments []transcriptionSegment `json:"segments"`
	XGroq    struct {
		ID string `json:"id"`
	} `json:"x_groq"`
}
