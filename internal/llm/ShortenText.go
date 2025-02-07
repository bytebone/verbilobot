package llm

import (
	"context"

	"github.com/conneroisu/groq-go"
)

// Removes speech patterns and filler words, returning the shortest possible form of the text without removing or altering any of the meaningful content.
func ShortenText(ctx context.Context, c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.ChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama318BInstant,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.RoleSystem,
				Content: "# TASK\nYou will be provided a voice transcript. Remove any speech patterns and filler words that are exclusive to spoken text, reducing the text to its shortest possible form without removing or changing any of the relevant contents. \n\n# RULES\n- Do not add formatting or markdown to the output. \n- Do not remove full sentences from the input UNLESS they add nothing meaningfull to the message\n- Separate different topics using TWO line breaks\n- Maintain the input language in your output.",
			},
			{
				Role:    groq.RoleUser,
				Content: text,
			},
		},
	})

	shortText = resp.Choices[0].Message.Content

	return
}
