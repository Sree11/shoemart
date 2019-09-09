package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch"
)

func regularQuery(es *SearchRequest) string {
	var regularQ = `
	{
		"query": {
			"multi_match": {
				"query": "` + strings.ToLower(es.Query) + `",
                "fields":["name", "brand", "colors", "currency"]
				}
			},
			"sort":{"` + strings.ToLower(es.SortKey) + `":"` + strings.ToLower(es.SortValue) + `"},
			"size":"` + fmt.Sprintf("%d", es.Limit) + `",
			"from":"` + strconv.Itoa(es.Offset) + `"
				
	}`
	return regularQ

}
func filterQuery(es *SearchRequest) string {
	var filterQ = ` {
		"query": {
		  "bool": {
			"must": {
			  "match": {
				"name": "` + strings.ToLower(es.Query) + `"
			  }
			},
			"filter":{
			"match":{
			"` + strings.ToLower(es.FilterKey) + `":"` + strings.ToLower(es.FilterValue) + `"}}
		  }
		},
		"sort":{"` + strings.ToLower(es.SortKey) + `":"` + strings.ToLower(es.SortValue) + `"},
			"size":"` + strconv.Itoa(es.Limit) + `",
			"from":"` + strconv.Itoa(es.Offset) + `"
	  }`
	return filterQ

}

func esQuery(es *SearchRequest) *strings.Reader {
	var query string
	if es.NoFilter {
		query = regularQuery(es)
	} else {
		query = filterQuery(es)
	}
	log.Println("\nquery:", query)

	isValid := json.Valid([]byte(query))

	if isValid == false {
		log.Println("Contruct error query not valid", query)
		log.Println("using default match_all")
		query = "{}"
	} else {
		log.Println("constructQuery is valid json", isValid)
	}

	var b strings.Builder
	b.WriteString(query)

	read := strings.NewReader(b.String())
	return read

}

func esGet(read *strings.Reader, c chan<- SearchMetrics) {
	var wg sync.WaitGroup

	log.SetFlags(0)
	ctx := context.Background()

	esCfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
		return
	}

	var mapResp map[string]interface{}
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("encoding error %s", err)
		return
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex("shoemart"),
		client.Search.WithBody(read),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithHuman(),
	)

	if err != nil {
		log.Fatalf("Elastic search error %s ", err)

	} else {
		log.Println("Response Type is ", reflect.TypeOf(res))
		defer res.Body.Close()
	}

	if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		log.Fatalf("Fatal error parsing json %s", err)
		return
	}

	esMetrics := SearchMetrics{}
	esMetrics.Hits = int(mapResp["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	esMetrics.Response = res.StatusCode
	esMetrics.Time = int(mapResp["took"].(float64))

	log.Printf("[%s] %d hits; took: %dms", res.Status(), int(mapResp["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(mapResp["took"].(float64)))

	for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {

		esMetrics.SearchResults = append(esMetrics.SearchResults, SearchResult{
			Brand:      string(hit.(map[string]interface{})["_source"].(map[string]interface{})["brand"].(string)),
			Name:       string(hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string)),
			Currency:   string(hit.(map[string]interface{})["_source"].(map[string]interface{})["currency"].(string)),
			Price:      float64(hit.(map[string]interface{})["_source"].(map[string]interface{})["price"].(float64)),
			Color:      string(hit.(map[string]interface{})["_source"].(map[string]interface{})["colors"].(string)),
			Sizes:      string(hit.(map[string]interface{})["_source"].(map[string]interface{})["sizes"].(string)),
			Categories: string(hit.(map[string]interface{})["_source"].(map[string]interface{})["categories"].(string)),
		})

	}
	c <- esMetrics
	defer wg.Wait()
	return

}
