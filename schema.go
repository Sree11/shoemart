package main

//SearchRequest will store the Request parameters
type SearchRequest struct {
	Query       string
	FilterKey   string
	FilterValue string
	NoFilter    bool
	SortKey     string
	SortValue   string
	Offset      int
	Limit       int
}

//SearchResult will store the Search results from Elastic Search
type SearchResult struct {
	Brand      string  `json:"brand"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Currency   string  `json:"currency"`
	Color      string  `json:"colors"`
	Sizes      string  `json:"sizes"`
	Categories string  `json:"categories"`
}

//SearchMetrics for Search query metrics returned by Elastic Search
type SearchMetrics struct {
	Response      int
	Hits          int
	Time          int
	SearchResults []SearchResult
}
