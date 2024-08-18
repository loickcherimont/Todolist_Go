package main

import (
	"database/sql" // required for db
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql" // required for mysql
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

// ******* STRUCTS *******
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

// ******* GLOBAL *******
var db *sql.DB

// Sample data
// todo: to insert into DB
var data = TodoPageData{
	Todos: []Todo{
		{Id: 1, Title: "Code a new project", Done: true},
		{Id: 2, Title: "Cook a healthier meal", Done: false},
		{Id: 3, Title: "Ride more often the bike", Done: false},
	},
	Username: "AdminDev",
}

// Generate a random key for cookie
var key = []byte(securecookie.GenerateRandomKey(32))
var store = sessions.NewCookieStore(key)

const PORT string = "3000"
const CRYPTCOST int = 14

// ******* MAIN *******
func main() {
	connectTo("todolist", "mysql")
	serveStatic()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/app", appHandler)

	fmt.Printf("\nServer OK and listening on http://localhost:%s/login \nTo stop it press, Ctrl+C\n**************************************************\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

// ******* FUNCTIONS *******

// Connect to a database system
// With the database given in parameter
func connectTo(dbName, dbSystem string) {

	var err error

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: dbName,
	}

	db, err = sql.Open(dbSystem, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("\n**************************************************\n\nConnection '%s':\n- Database: %s\n- Result: SUCCESS\n", dbSystem, dbName)
}

// Handle static files
func serveStatic() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

// Return an error if input is empty
func handleError(input string) ErrorItem {
	if len(input) == 0 {
		return ErrorItem{IsError: true, ErrorMessage: fmt.Errorf("Sorry, I cannot add empty tasks!")}
	}

	return ErrorItem{IsError: false, ErrorMessage: nil}
}

// Revoke users authentication
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Close session
	session.Values["authenticated"] = false
	
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")

	tmpl := template.Must(template.ParseFiles("view.html"))

	// Prevent user try to access without authentication
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusForbidden)
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

	tmpl.ExecuteTemplate(w, "view.html", data)
	return
}

// Authentify user
// Or just execute login page template
func loginHandler(w http.ResponseWriter, r *http.Request) {

	var u User // useful?

	tmpl := template.Must(template.ParseGlob("*.html"))

	// User connect to the app
	// todo: separate form handler from login page handler (possible)
	if r.Method == http.MethodPost {

		// Authentication
		username := r.PostFormValue("login")
		password := r.PostFormValue("password")

		// Retrieve data from DB
		row := db.QueryRow("SELECT * FROM users WHERE username = ? AND password = ?", username, password)

		// If error is found, show it using UI message
		if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
			fmt.Errorf("Connection error : %v", err)
			data.Errors.ErrorMessage = fmt.Errorf("Login or password is wrong! Retry please!")
			data.Errors.IsError = true
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Start new session
		session, _ := store.Get(r, "cookie-name")

		// Set user as authenticated
		// Remove previous errors
		// And execute appHandler
		session.Values["authenticated"] = true
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Errors.IsError = false
		data.Errors.ErrorMessage = nil
		http.Redirect(w, r, "/app", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "login.html", data)
	return
}

