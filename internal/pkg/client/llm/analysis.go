package llm

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
)

// SentimentAnalysisResponse response from sentiment analysis
type SentimentAnalysisResponse struct {
	Sentiment int `json:"sentiment"`
}

// SentimentAnalysisRequest request for sentiment analysis
type SentimentAnalysisRequest struct {
	Context string `json:"context"`
}

// SentimentAnalysisRequest validate sentiment analysis request
func (sa *SentimentAnalysisRequest) Validate() error {
	var errs []error

	if strings.TrimSpace(sa.Context) == "" {
		errs = append(errs, errors.New("context is required"))
	}

	return errors.Join(errs...)
}

// SentimentAnalysis perform sentiment analysis
func (c *Client) SentimentAnalysis(ctx context.Context, req *SentimentAnalysisRequest) (*SentimentAnalysisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("wrong request: %w", err)
	}

	// TODO: раскомментировать после реализации LLM сервиса
	// llmResp, err := c.makeRequestJSON(ctx, http.MethodPost, c.endpoints.SentimentAnalysis, req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to send request to %s: %w", c.endpoints.SentimentAnalysis, err)
	// }

	// var resp SentimentAnalysisResponse
	// if err := json.Unmarshal(llmResp, &resp); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	// }

	// TODO: удалить заглушку после реализации LLM сервиса
	var resp SentimentAnalysisResponse
	resp.Sentiment = int((2*rand.Float32() - 1.0) * 100)
	// TODO: удалить заглушку после реализации LLM сервиса

	return &resp, nil
}
