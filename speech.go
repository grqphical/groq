package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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

// Transcribes a given audio file using one of Groq's hosted Whipser models
func (g *GroqClient) TranscribeAudio(filename string, model string, config *TranscriptionConfig) (Transcription, error) {
	req, err := createGroqRequest("audio/transcriptions", g.apiKey, "POST", nil)
	if err != nil {
		return Transcription{}, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return Transcription{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return Transcription{}, err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return Transcription{}, err
	}

	err = writer.WriteField("model", model)
	if err != nil {
		return Transcription{}, err
	}

	if config != nil {
		if config.Language != "" {
			err = writer.WriteField("language", config.Language)
			if err != nil {
				return Transcription{}, err
			}
		}
		if config.Prompt != "" {
			err = writer.WriteField("prompt", config.Prompt)
			if err != nil {
				return Transcription{}, err
			}
		}
		if config.ResponseFormat != "" {
			err = writer.WriteField("response_format", config.ResponseFormat)
			if err != nil {
				return Transcription{}, err
			}
		}
		if config.Temperature != 0 {
			err = writer.WriteField("temperature", fmt.Sprintf("%f", config.Temperature))
			if err != nil {
				return Transcription{}, err
			}
		}
	}

	err = writer.Close()
	if err != nil && err != io.EOF {
		return Transcription{}, err
	}

	req.Body = io.NopCloser(body)

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Transcription{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Transcription{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var responseTranscription Transcription

	if config != nil && config.ResponseFormat == "text" {
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			return Transcription{}, err
		}
		responseTranscription.Text = string(text)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&responseTranscription)
		if err != nil {
			return Transcription{}, err
		}
	}

	return responseTranscription, nil
}

// Translates a given audio file into English.
func (g *GroqClient) TranslateAudio(filename string, model string, config *TranslationConfig) (Transcription, error) {
	req, err := createGroqRequest("audio/translations", g.apiKey, "POST", nil)
	if err != nil {
		return Transcription{}, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return Transcription{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return Transcription{}, err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return Transcription{}, err
	}

	err = writer.WriteField("model", model)
	if err != nil {
		return Transcription{}, err
	}

	if config != nil {
		if config.Prompt != "" {
			err = writer.WriteField("prompt", config.Prompt)
			if err != nil {
				return Transcription{}, err
			}
		}
		if config.ResponseFormat != "" {
			err = writer.WriteField("response_format", config.ResponseFormat)
			if err != nil {
				return Transcription{}, err
			}
		}
		if config.Temperature != 0 {
			err = writer.WriteField("temperature", fmt.Sprintf("%f", config.Temperature))
			if err != nil {
				return Transcription{}, err
			}
		}
	}

	err = writer.Close()
	if err != nil && err != io.EOF {
		return Transcription{}, err
	}

	req.Body = io.NopCloser(body)

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Transcription{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Transcription{}, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var responseTranslation Transcription

	if config != nil && config.ResponseFormat == "text" {
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			return Transcription{}, err
		}
		responseTranslation.Text = string(text)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&responseTranslation)
		if err != nil {
			return Transcription{}, err
		}
	}

	return responseTranslation, nil
}
