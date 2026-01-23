package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type chatReq struct {
	Model    string     `json:"model"`
	Messages []chatMsg  `json:"messages"`
	
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResp struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateMarkdownFromTxt(inputPath, outputPath string) error {

	err := godotenv.Load(".env")
	if err != nil  {
		log.Println(err)
	}

	apiKey := os.Getenv("apiKey") // api key

	if apiKey == "" {
		return errors.New("DEEPSEEK_API_KEY is empty (set env var)")
	}

	in, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input txt: %w", err)
	}

	reqBody := chatReq{
		Model: "deepseek-chat",
		Messages: []chatMsg{
			{
				Role: "user",
				Content: "Return a single Markdown document.\n\n" + string(in),
			},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(b))
	// httpReq, err := http.NewRequest("POST", "https://api.artemox.com/v1", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request deepseek: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("deepseek status %d: %s", resp.StatusCode, string(raw))
	}

	var out chatResp
	if err := json.Unmarshal(raw, &out); err != nil {
		return fmt.Errorf("unmarshal response: %w; raw=%s", err, string(raw))
	}

	if len(out.Choices) == 0 {
		return errors.New("deepseek: empty choices")
	}

	md := out.Choices[0].Message.Content
	if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
		return fmt.Errorf("write output md: %w", err)
	}

	return nil
}
