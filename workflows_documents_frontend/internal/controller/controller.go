package controller

import (
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println(username, password)
	if true {
		w.Header().Set("HX-Redirect", "/home")
		return
	}
}
