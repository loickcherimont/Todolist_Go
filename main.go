package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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
	Todos    []Todo
	Errors   ErrorItem
	Username string
}

type User struct {
	ID       int64
	Username string
	Password string
}

// Some data
// Before DB
var data = TodoPageData{
	Todos: []Todo{
		{Id: 1, Title: "Code a new project", Done: true},
		{Id: 2, Title: "Cook a healthier meal", Done: false},
		{Id: 3, Title: "Ride more often the bike", Done: false},
	},
	Username: "AdminDev",
}

var isConnected bool

const PORT string = "3000"

// Return an error if input is empty
func handleError(input string) ErrorItem {
	if len(input) == 0 {
		return ErrorItem{IsError: true, ErrorMessage: fmt.Errorf("Sorry, I cannot add empty tasks!")}
	}

	return ErrorItem{IsError: false, ErrorMessage: nil}
}

// Authentication system before DB
// Ancient indexHandler
func loginHandler(w http.ResponseWriter, r *http.Request) {

	var u User

	tmpl := template.Must(template.ParseGlob("*.html"))

	if r.Method == http.MethodPost {
		username := r.PostFormValue("login")
		password := r.PostFormValue("password")

		// Retrieve data from DB
		row := db.QueryRow("SELECT * FROM users WHERE username = ? AND password = ?", username, password)

		// If error is found show error message to user
		if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
			if err == sql.ErrNoRows {
				fmt.Errorf("Connection error : %v", err)
				data.Errors.ErrorMessage = fmt.Errorf("Login or password is wrong! Retry please!")
				data.Errors.IsError = true
			}
			fmt.Errorf("Connection error : %v", err)
			data.Errors.ErrorMessage = fmt.Errorf("Login or password is wrong! Retry please!")
			data.Errors.IsError = true
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Delete previous errors
		// From login step
		isConnected = true
		data.Errors.IsError = false
		data.Errors.ErrorMessage = nil
		http.Redirect(w, r, "/app", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "login.html", data)
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	isConnected = false // Change user connection status
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("view.html"))

	// Prevent user try to access without authentication
	if !isConnected {
		fmt.Fprintln(w, "<h1>Sorry but access denied!</h1>")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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

	// todo: complete this function using the "ancient indexHandler" block
	tmpl.ExecuteTemplate(w, "view.html", data)
	return
}

// Serve CSS file
func serveStatic() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	connectTo("todolist")
	serveStatic()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/app", appHandler)
	fmt.Printf("\nServer OK and listening on http://localhost:%s/login\n To stop it press, Ctrl+C\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func connectTo(dbName string) {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: dbName,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

// Last change : We want to establish a login/logout using DB
