package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//Users holds values
var Users = map[string]string{
	"username": "admin",
	"password": "user",
	"hash":     "$2y$14$BrWf3t4LSltZDMDOVM8/guQpmuFnR731rKMY27H77aJZ67Iy3n2nq",
}

func enforceAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	enteredpassword := r.PostForm.Get("password")
	password := Users["hash"]
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(enteredpassword))
	fmt.Println(err)
	if err != nil {
		fmt.Printf("Incorrect Password")
		http.Redirect(w, r, "/", 302)
	}
	log.Printf("User Authenticated")
	session, _ := store.Get(r, "session")

	session.Values["username"] = username
	session.Options.MaxAge = 60 * 5
	session.Save(r, w)
	http.Redirect(w, r, "/v1/products", 302)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "login.html", nil)

}
