package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	Id    string
	Title string
	Done  bool
}

type TodoPageData struct {
	Todos []Todo
}

var data = TodoPageData{
	Todos: []Todo{
		{Id: "T1", Title: "Code a new project", Done: true},
		{Id: "T2", Title: "Cook a healthier meal", Done: false},
		{Id: "T3", Title: "Ride more often the bike", Done: false},
	},
}

// Parse and execute the main template
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view.html"))
	tmpl.Execute(w, data)
}

// Serve CSS file
func serveStatic() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	const PORT string = "3000"
	serveStatic()
	http.HandleFunc("/", indexHandler)
	fmt.Printf("\nServer OK and listening on http://localhost:%s\n To stop it press, Ctrl+C", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
