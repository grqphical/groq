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

// Transcribes a given audio file using one of Groq's hosted Whipser models
func (g *GroqClient) TranscribeAudio(filename string, model string, config *TranscriptionConfig) (Transcription, error) {
	req, err := createGroqRequest("audio/transcriptions", g.apiKey, "POST")
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
