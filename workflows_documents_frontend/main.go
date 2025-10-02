package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/rcarvalho-pb/workflows-document_frontend/internal/controller"
	"github.com/rcarvalho-pb/workflows-document_frontend/internal/pages"
)

const PORT = 8080

func main() {
	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.Handle("/", templ.Handler(pages.LoginPage()))
	r.Handle("/home", templ.Handler(pages.HomePage()))
	r.HandleFunc("/login-func", controller.Login)

	log.Printf("server starting on port :%d\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), r))
}
