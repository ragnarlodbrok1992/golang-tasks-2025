package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/json", jsonHandler)

	// Start the server
	port := ":8080"
	log.Printf("Server is listening on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// rootHandler handles requests to the root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, 200)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the root path!")
}

// helloHandler handles requests to /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, 200)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!")
}

// jsonHandler handles requests to /json and returns a JSON response
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, 200)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, JSON!"})
}

// logRequest logs the method, path, and status code of each request
func logRequest(r *http.Request, statusCode int) {
	log.Printf("%s %s %d", r.Method, r.URL.Path, statusCode)
}
