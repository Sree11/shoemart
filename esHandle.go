package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func productGetHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "product.html", nil)

}

//SearchHandler will serve the searches
func productPostHandler(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	var es SearchRequest

	r.ParseForm()
	user := r.PostForm.Get("username")
	fmt.Println(user)

	val := r.PostForm.Get("searchquery")
	fmt.Println(val)
	if val != "" {
		es.Query = val
	} else {
		log.Printf("no search")
		http.Redirect(w, r, "/products", 302)

	}

	val = r.PostForm.Get("filter")
	if val != "" {
		filterArray := strings.Split(val, ":")
		es.FilterKey = filterArray[0]
		es.FilterValue = filterArray[1]
	} else {
		es.NoFilter = true
	}

	val = r.PostForm.Get("sort")
	if val != "" {
		sortArray := strings.Split(val, ":")
		es.SortKey = sortArray[0]
		es.SortValue = sortArray[1]
	} else {
		//if Sort is not given default to sort on Score
		es.SortKey = "_score"
		es.SortValue = "desc"
	}

	val = r.PostForm.Get("offset")
	if val != "" {
		es.Offset, _ = strconv.Atoi(val)
	} else {
		//Set default offset to 1
		es.Offset = 0
	}
	val = r.PostForm.Get("limit")
	if val != "" {
		es.Limit, _ = strconv.Atoi(val)
	} else {
		//set default limit to the shard limit
		es.Limit = 2000
	}
	templates.ExecuteTemplate(w, "productresult.html", nil)
	defer r.Body.Close()
	//qr := constructQuery(&es)
	//log.Println(qr)
	//rd := esQuery(qr)
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

	defer wg.Done()
	return

}
