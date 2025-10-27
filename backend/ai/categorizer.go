package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

type CategorySuggestion struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// CategorizeBookmark uses Claude API to suggest categories for a bookmark
func CategorizeBookmark(ctx context.Context, tweetText string, existingCategories []string) ([]string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	// Build the prompt
	prompt := buildCategorizationPrompt(tweetText, existingCategories)

	// Create the request
	reqBody := ClaudeRequest{
		Model:     "claude-3-haiku-20240307", // Cheapest model
		MaxTokens: 200,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API call
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	// Extract categories from the response
	categories := parseCategories(claudeResp.Content[0].Text)
	return categories, nil
}

func buildCategorizationPrompt(tweetText string, existingCategories []string) string {
	categoriesStr := "None yet"
	if len(existingCategories) > 0 {
		categoriesStr = strings.Join(existingCategories, ", ")
	}

	return fmt.Sprintf(`You are a bookmark categorization assistant. Analyze this tweet and suggest 1-2 relevant categories.

Tweet text: "%s"

Existing user categories: %s

Instructions:
1. If the tweet matches existing categories, use those (exact name match)
2. If no match, suggest NEW categories (max 2)
3. Use simple, clear category names (e.g., "Tech", "Coding", "Design", "AI", "Business", "Marketing")
4. Respond with ONLY category names, comma-separated
5. No explanations, just category names

Example responses:
- "Tech, AI"
- "Design"
- "Marketing, Business"

Your response:`, tweetText, categoriesStr)
}

func parseCategories(response string) []string {
	// Clean up the response
	response = strings.TrimSpace(response)
	response = strings.Trim(response, "\"'")
	
	// Split by comma
	parts := strings.Split(response, ",")
	
	categories := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && len(part) < 50 { // Sanity check
			categories = append(categories, part)
		}
	}
	
	// Limit to 2 categories
	if len(categories) > 2 {
		categories = categories[:2]
	}
	
	return categories
}

// BatchCategorizeBookmarks categorizes multiple bookmarks in a single operation
func BatchCategorizeBookmarks(ctx context.Context, bookmarks []struct {
	ID        string
	TweetText string
}, existingCategories []string) (map[string][]string, error) {
	results := make(map[string][]string)
	
	for _, bookmark := range bookmarks {
		categories, err := CategorizeBookmark(ctx, bookmark.TweetText, existingCategories)
		if err != nil {
			// Log error but continue with other bookmarks
			fmt.Printf("Error categorizing bookmark %s: %v\n", bookmark.ID, err)
			continue
		}
		results[bookmark.ID] = categories
	}
	
	return results, nil
}

