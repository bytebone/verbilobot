package fileutils

import (
	"context"
	"os"

	"github.com/conneroisu/groq-go"
)

func Transcribe(ctx context.Context, file *os.File) (text string, err error) {
	groqClient, err := groq.NewClient(os.Getenv("VERBILO_GROQ_TOKEN"))
	if err != nil {
		return "", err
	}
	resp, err := groqClient.Transcribe(ctx, groq.AudioRequest{
		Model:    groq.ModelWhisperLargeV3Turbo,
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
