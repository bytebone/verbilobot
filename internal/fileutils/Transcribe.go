package fileutils

import (
	"context"
	"os"

	"github.com/conneroisu/groq-go"
)

/*
func transcribe(file *os.File) (text string, err error) {
	apiKey := os.Getenv("VERBILO_GROQ_TOKEN")

	// f, err := os.Open(file.Name())
	// if err != nil {
	// 	return
	// }
	// defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return
	}

	err = writer.WriteField("model", "whisper-large-v3")
	if err != nil {
		return
	}
	err = writer.WriteField("response_format", "text")
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/audio/transcriptions", &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// req.Body = io.NopCloser(buf)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	} else if resp.StatusCode == 429 {
		return "", fmt.Errorf(resp.Status)
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(respBody), nil
}
*/

func Transcribe(ctx context.Context, file *os.File) (text string, err error) {
	groqClient, err := groq.NewClient(os.Getenv("VERBILO_GROQ_TOKEN"))
	if err != nil {
		return "", err
	}
	resp, err := groqClient.Transcribe(ctx, groq.AudioRequest{
		Model:    groq.WhisperLargeV3,
		Format:   groq.FormatText,
		FilePath: file.Name(),
	})
	if err != nil {
		return "", err
	}
	// if resp == "" {
	// 	return "", fmt.Errorf("empty response")
	// }

	text = resp.Text

	return
}
