package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var templates *template.Template
var store = sessions.NewCookieStore([]byte("secret"))

//EsInit for init
func EsInit() (esClient *elasticsearch.Client) {
	log.SetFlags(0)

	esCfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
		return
	}
	return client
}

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	r := mux.NewRouter()

	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", loginPostHandler).Methods("POST")

	AddV1Routes(r.PathPrefix("/v1").Subrouter())

	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	//r.PathPrefix("/v1/products/").Handler(http.StripPrefix("/v1/products/", fs))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))

}

//AddV1Routes adds version to URL
func AddV1Routes(r *mux.Router) {
	r.HandleFunc("/products", enforceAuth(productGetHandler)).Methods("GET")
	r.HandleFunc("/products", enforceAuth(productPostHandler)).Methods("POST")

	r.HandleFunc("/search", enforceAuth(searchGetHandler)).Methods("GET")

}
