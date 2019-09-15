package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func regularQuery(es *SearchRequest) string {
	var rq RegularQuery
	fieldstr := []string{"name", "brand", "colors", "currency"}
	SortField := map[string]interface{}{"price": "desc"}
	for _, item := range fieldstr {
		rq.Query.MultiMatch.Fields = append(rq.Query.MultiMatch.Fields, item)
	}
	rq.Query.MultiMatch.Query = strings.ToLower(es.Query)
	rq.Sort = SortField
	rq.Size = strconv.Itoa(es.Limit)
	rq.From = strconv.Itoa(es.Offset)
	jsonrq, _ := rq.Marshal()
	return string(jsonrq)

}
func filterQuery(es *SearchRequest) string {
	var fq FilterQuery
	filterstr := map[string]interface{}{strings.ToLower(es.FilterKey): strings.ToLower(es.FilterValue)}
	SortField := map[string]interface{}{strings.ToLower(es.SortKey): strings.ToLower(es.SortValue)}

	fq.FQuery.Bool.Must.Match.Name = strings.ToLower(es.Query)
	fq.FQuery.Bool.Filter.Match = filterstr

	fq.Sort = SortField
	fq.Size = strconv.Itoa(es.Limit)
	fq.From = strconv.Itoa(es.Offset)
	jsonfq, _ := fq.Marshal()
	fmt.Printf(string(jsonfq))
	return string(jsonfq)

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

	var mapResp map[string]interface{}
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("encoding error %s", err)
		return
	}

	ctx := context.Background()
	client := EsInit()
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
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	esresults, eserr := UnmarshalESResult(body)
	if eserr != nil {
		fmt.Println("error parsing with unmarshall ", eserr)
	}
	fmt.Println("Here it comes ========")
	fmt.Println(esresults)

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
