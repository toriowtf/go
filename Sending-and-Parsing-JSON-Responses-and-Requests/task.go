package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestBody is a struct to represent the expected JSON data structure.
type RequestBody struct {
	Message string `json:"message"`
}

// ResponseBody is a struct to represent the JSON response structure.
type ResponseBody struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", handleRequest)
	port := 8080
	fmt.Printf("Server is listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON data from the request body
	var requestBody RequestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Check if the "message" field is absent or empty
	if requestBody.Message == "" {
		errorResponse := ResponseBody{
			Status:  "400",
			Message: "Invalid JSON message",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Print the received message to the server console
	fmt.Printf("Received message: %s\n", requestBody.Message)

	// Send a JSON response with a success message
	response := ResponseBody{
		Status:  "success",
		Message: "Data successfully received",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}