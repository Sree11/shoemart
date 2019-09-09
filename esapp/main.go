package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var templates *template.Template
var store = sessions.NewCookieStore([]byte("secret"))

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
	//r.HandleFunc("/search", enforceAuth(searchPostHandler)).Methods("POST")

}
