package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ErrorItem struct {
	IsError      bool
	ErrorMessage error
}

type Todo struct {
	Id    int
	Title string
	Done  bool
}

type TodoPageData struct {
	Todos  []Todo
	Errors ErrorItem
}

// Some data
// Before DB
var data = TodoPageData{
	Todos: []Todo{
		{Id: 1, Title: "Code a new project", Done: true},
		{Id: 2, Title: "Cook a healthier meal", Done: false},
		{Id: 3, Title: "Ride more often the bike", Done: false},
	},
}

var err ErrorItem

// Return an error if input is empty
func handleError(input string) ErrorItem {
	if len(input) == 0 {
		return ErrorItem{IsError: true, ErrorMessage: fmt.Errorf("Sorry, I cannot add empty tasks!")}
	}

	return ErrorItem{IsError: false, ErrorMessage: nil}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Parse all html files
	tmpl := template.Must(template.ParseGlob("*.html"))

	// Process POSTed data
	if r.Method == http.MethodPost {

		// Process errors
		data.Errors.IsError = handleError(r.PostFormValue("input_task")).IsError
		data.Errors.ErrorMessage = handleError(r.PostFormValue("input_task")).ErrorMessage

		// Prevent task adding
		// If error detected
		if !data.Errors.IsError {
			data.Todos = append(data.Todos, Todo{Id: len(data.Todos) + 1, Title: r.PostFormValue("input_task"), Done: false})
		}
	}

	tmpl.Execute(w, data)
	return
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
