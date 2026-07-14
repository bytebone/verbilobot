package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type TeableResult struct {
	Status         string
	OriginalSize   float64
	CompressedSize float64
	Duration       float64
	Err            string
	ErrorStage     string
	FileType       string
}

func Teable(_ context.Context, r TeableResult) {
	host := os.Getenv("VERBILO_TEABLE_HOST")
	tableID := os.Getenv("VERBILO_TEABLE_TABLE_ID")
	apiKey := os.Getenv("VERBILO_TEABLE_API_KEY")
	if host == "" || tableID == "" || apiKey == "" {
		return
	}

	payload := map[string]any{
		"fieldKeyType": "dbFieldName",
		"records": []map[string]any{
			{
				"fields": map[string]any{
					"status":          r.Status,
					"original_size":   r.OriginalSize,
					"compressed_size": r.CompressedSize,
					"duration":        r.Duration,
					"error":           r.Err,
					"error_stage":     r.ErrorStage,
					"filetype":        r.FileType,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("teable: marshal: %v", err)
		return
	}

	url := strings.TrimRight(host, "/") + "/api/table/" + tableID + "/record"

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			log.Printf("teable: new request: %v", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("teable: post: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			log.Printf("teable: unexpected status %d: %s", resp.StatusCode, string(b))
		}
	}()
}
