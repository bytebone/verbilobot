package llm

import (
	"github.com/jpoz/groq"
)

// Removes speech patterns and filler words, returning the shortest possible form of the text without removing or altering any of the meaningful content.
func ShortenText(c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.CreateChatCompletion(groq.CompletionCreateParams{
		Model: "llama-3.1-8b-instant",
		Messages: []groq.Message{
			{
				Role:    "system",
				Content: "You will be provided a voice transcript. Your task is to remove any speech patterns and filler words that are exclusive to spoken text, and thereby reduce the text to its shortest possible form without removing or altering any of the meaningful contents whatsoever. Do not add formatting, markdown or additional text of your own around the shortened text. The text may be provided in any language, make sure that your shortened version matches the input language.",
			},
			{
				Role:    "user",
				Content: text,
			},
		},
	})

	shortText = resp.Choices[0].Message.Content

	return
}
