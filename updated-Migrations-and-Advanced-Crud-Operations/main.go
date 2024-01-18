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
var dbName = "registration"
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
	r.HandleFunc("/login", LoginHandler).Methods("GET")
	r.HandleFunc("/login", LoginPostHandler).Methods("POST")
	r.HandleFunc("/forgot-password", ForgotPasswordHandler).Methods("GET")
	r.HandleFunc("/forgot-password", ForgotPasswordPostHandler).Methods("POST")	

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Check user credentials in MongoDB (assuming MongoDB schema is as described in previous responses)
	if checkCredentials(username, password) {
		// Successful login, redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Failed login, render login page with an error message
		renderTemplate(w, "login", map[string]interface{}{"Error": "Invalid username or password"})
	}
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "forgot-password", nil)
}

func ForgotPasswordPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get email from the form
	user := r.Form.Get("username")
	password := r.Form.Get("password")

	// Reset password logic (assuming MongoDB schema is as described in previous responses)
	if resetPassword(user, password) {
		// Password reset successful, redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// Password reset failed, render forgot password page with an error message
		renderTemplate(w, "forgot-password", map[string]interface{}{"Error": "User does not exist."})
	}
}

// Check user credentials in MongoDB
func checkCredentials(username, password string) bool {
	collection := client.Database(dbName).Collection(collectionName)

	filter := map[string]interface{}{
		"username": username,
		"password": password,
	}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Println("Error checking credentials:", err)
		return false
	}

	return count > 0
}

// Reset password in MongoDB
func resetPassword(username, password string) bool {
	collection := client.Database(dbName).Collection(collectionName)

	// Check if the email exists in the database
	filter := map[string]interface{}{
		"username": username,
	}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Println("Error checking user:", err)
		return false
	}

	if count > 0 {
		// Perform password reset logic here (e.g., generate a new password and update the database)
		// For simplicity, let's assume a hardcoded new password "newpassword" for this example
		newPassword := password

		update := map[string]interface{}{
			"$set": map[string]interface{}{
				"password": newPassword,
			},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			fmt.Println("Error updating password:", err)
			return false
		}

		return true
	}

	return false
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