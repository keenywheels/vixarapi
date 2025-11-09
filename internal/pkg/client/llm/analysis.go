package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SentimentAnalysisResponse response from sentiment analysis
type SentimentAnalysisResponse struct {
	Sentiment int16 `json:"sentiment"`
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

	llmResp, err := c.makeRequestJSON(ctx, http.MethodPost, c.endpoints.SentimentAnalysis, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", c.endpoints.SentimentAnalysis, err)
	}

	var resp SentimentAnalysisResponse
	if err := json.Unmarshal(llmResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &resp, nil
}
