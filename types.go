package groq

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
	// What language the audio is in. if blank the model will guess it
	Language string
	// An optional text to guide the model's style or continue a previous audio segment. The prompt should match the audio language.
	Prompt string
	// The format of the transcript output, in one of these options: json, text, or verbose_json
	ResponseFormat string
	// The sampling temperature, between 0 and 1.
	Temperature float64
}

// TranslationConfig houses configuration options for translation requests
type TranslationConfig struct {
	// An optional text to guide the model's style or continue a previous audio segment. The prompt should match the audio language.
	Prompt string
	// The format of the transcript output, in one of these options: json, text, or verbose_json
	ResponseFormat string
	// The sampling temperature, between 0 and 1.
	Temperature float64
}

// Transcription represents an audio transcription/translation result from one of Groq's models
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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Conversation is a struct that allows you to construct chat completion requests
type Conversation struct {
	messages     []Message
	systemPrompt string
}

type ChatCompletionRequest struct {
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// Maximum amount of tokens that can be generated in the completion
	MaxTokens int `json:"max_tokens,omitempty"`

	Messages []Message `json:"messages"`

	Model string `json:"model"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	ResponseFormat struct {
		Type string `json:"type,omitempty"`
	} `json:"response_format,omitempty"`

	Seed   int      `json:"seed,omitempty"`
	Stop   []string `json:"stop,omitempty"`
	Stream bool     `json:"stream,omitempty"`
	// The sampling temperature, between 0 and 1.
	Temperature float64 `json:"temperature,omitempty"`
	User        string  `json:"user,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
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

	// Whether or not the API should stream responses
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
