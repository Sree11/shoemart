package main

import "encoding/json"

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

// type Sort struct {
// 	Score string `json:"_score"`
// }
type MultiMatch struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}
type Query struct {
	MultiMatch MultiMatch `json:"multi_match"`
}
type RegularQuery struct {
	Query Query                  `json:"query"`
	Sort  map[string]interface{} `json:"sort"`
	Size  string                 `json:"size"`
	From  string                 `json:"from"`
}

//Marshal the unmarshall
func (r *RegularQuery) Marshal() ([]byte, error) {
	return json.MarshalIndent(r, "", " ")
}

//==================================================

type FilterQuery struct {
	FQuery FQuery                 `json:"query"`
	Sort   map[string]interface{} `json:"sort"`
	Size   string                 `json:"size"`
	From   string                 `json:"from"`
}

type FQuery struct {
	Bool Bool `json:"bool"`
}

type Bool struct {
	Must   Must   `json:"must"`
	Filter Filter `json:"filter"`
}

type Filter struct {
	Match map[string]interface{} `json:"match"`
}

type Must struct {
	Match MustMatch `json:"match"`
}

type MustMatch struct {
	Name string `json:"name"`
}

func (r *FilterQuery) Marshal() ([]byte, error) {

	return json.MarshalIndent(r, "", " ")
}
