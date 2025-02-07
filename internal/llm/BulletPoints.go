package llm

import (
	"context"

	"github.com/conneroisu/groq-go"
)

// Rewrites the given text as bullet points, reducing the texts to the shortest possible length without cutting down on relevant information.
func BulletPoints(ctx context.Context, c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.ChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama318BInstant,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.RoleSystem,
				Content: "# TASK\nYou will be provided a voice transcript. Your task is to reduce the text to bullet points, reducing the texts to its core points.\n\n# RULES\n- Since the transcript may contain errors, if you find a sentence does not make sense, correct it at your own discretion for the output.\n- Do not add formatting or markdown to the output. \n- Keep formalities to a minimum and match the tone of the input.\n- Maintain the input language in your output.",
			},
			{
				Role:    groq.RoleUser,
				Content: text,
			},
			{
				Role:    groq.RoleAssistant,
				Content: "• ",
			},
		},
		Temperature: 0.65,
	})

	shortText = "• " + resp.Choices[0].Message.Content

	return
}
