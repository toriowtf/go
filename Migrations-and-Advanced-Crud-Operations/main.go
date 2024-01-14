// main.go
package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// MongoDB connection details
var mongoURI = "mongodb://localhost:27017"
var dbName = "userdb"
var collectionName = "users"

var client *mongo.Client

func init() {
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
}

func main() {
	r := mux.NewRouter()

	// Serve static files (e.g., CSS, JS)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Define routes
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/register", RegisterHandler).Methods("GET")
	r.HandleFunc("/register", RegisterPostHandler).Methods("POST")

	// Start the server
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handlers
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Insert data into MongoDB
	insertData(username, password)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Insert data into MongoDB
func insertData(username, password string) {
	collection := client.Database(dbName).Collection(collectionName)

	user := map[string]interface{}{
		"username": username,
		"password": password,
	}

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error inserting data into MongoDB:", err)
		return
	}

	fmt.Println("Data inserted into MongoDB successfully")
}
