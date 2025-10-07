package model

type APIBatchRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}
type APIBatchRequestA []APIBatchRequest

type APIBatchResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}
type APIBatchResponseA []APIBatchResponse

type APIShorURL struct {
	HTTPCode int    `json:"http_code,omitempty"`
	ShortURL string `json:"short_url,omitempty"`
}

type APIBatchShorURLs struct {
	HTTPCode int `json:"http_code,omitempty"`
	APIBatchResponseA
}

type APIRequest struct {
	URL string `json:"url,omitempty"`
}

type APIResult struct {
	Result string `json:"result,omitempty"`
}
