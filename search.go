package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func searchGetHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "search.html", nil)
	var wg sync.WaitGroup
	var es SearchRequest

	keys := r.URL.Query()
	fmt.Println(keys)
	if val, ok := keys["q"]; ok {
		es.Query = val[0]
	} else {
		log.Printf("redirect")
		log.Printf("no search")
		http.Redirect(w, r, "/search", 302)

	}

	if val, ok := keys["filter"]; ok {
		filterArray := strings.Split(val[0], ":")
		es.FilterKey = filterArray[0]
		es.FilterValue = filterArray[1]

	} else {
		es.NoFilter = true
	}

	if val, ok := keys["sort"]; ok {
		sortArray := strings.Split(val[0], ":")
		es.SortKey = sortArray[0]
		es.SortValue = sortArray[1]
	} else {
		//if Sort is not given default to sort on Score
		es.SortKey = "_score"
		es.SortValue = "desc"
	}

	if val, ok := keys["offset"]; ok {
		es.Offset, _ = strconv.Atoi(val[0])
	} else {
		//Set default offset to 1
		es.Offset = 0
	}
	if val, ok := keys["limit"]; ok {
		es.Limit, _ = strconv.Atoi(val[0])
	} else {
		//set default limit to the shard limit
		es.Limit = 2000
	}

	defer r.Body.Close()

	rd := esQuery(&es)

	c := make(chan SearchMetrics)
	wg.Add(1)
	go esGet(rd, c)

	metrics := <-c

	metricMessage := ` No of Hits ` + strconv.Itoa(metrics.Hits) + ` in ` + strconv.Itoa(metrics.Time) + ` msec`

	htmlbody := `                            
    <h6 class="text-center text-info">` + metricMessage + `</h6>
	<table class="table table-striped">
		<thead>
			<tr><th scope="col"></th><th scope="col">Brand</th><th scope="col"></th>
			<th scope="col">Name</th><th scope="col"></th><th scope="col">Currency</th>
			<th scope="col"></th><th scope="col">Price</th><th scope="col"></th>
			<th scope="col">Colors</th><th scope="col"></th><th scope="col">Sizes</th>
			</tr>
		</thead>`

	fmt.Fprintf(w, htmlbody)
	fmt.Fprintf(w, `<tbody>`)
	for _, val := range metrics.SearchResults {
		fmt.Fprintf(w, `<tr>`)
		fmt.Fprintf(w, `<th scope="row"></th>`)
		fmt.Fprintf(w, `<td>`+val.Brand+`<td>`)
		fmt.Fprintf(w, `<td>`+val.Name+`<td>`)
		fmt.Fprintf(w, `<td>`+val.Currency+`<td>`)
		fmt.Fprintf(w, `<td>`+strconv.FormatFloat(val.Price, 'f', 2, 64)+`<td>`)
		fmt.Fprintf(w, `<td>`+val.Color+`<td>`)
		fmt.Fprintf(w, `<td>`+val.Sizes+`<td>`)
		fmt.Fprintf(w, `</tr>`)

	}
	fmt.Fprintf(w, `<tbody></table>`)
	//templates.ExecuteTemplate(w, "productresults.html", nil)
	defer wg.Done()
	return

}
