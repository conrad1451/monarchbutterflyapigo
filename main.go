// --- main.go ---

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// The `pq` package is a pure Go PostgreSQL driver for `database/sql`.
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// MonarchRecord represents a row from the database table.
// Using a map[string]interface{} is flexible since the schema
// might change, similar to the dynamic dictionary creation in the Python version.
type MonarchRecord map[string]interface{}

// getMonarchsHandler fetches all Monarch butterfly data from the database
// and serves it as a JSON response.
func getMonarchsHandler(w http.ResponseWriter, r *http.Request) {
	// Use the same table name as the Python script for consistency.
	tableName := "june212025"
	 
	connStr := os.Getenv("GOOGLE_VM_DOCKER_HOSTED_SQL")
	
	// Open a connection to the database.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	// Ping the database to ensure the connection is live.
	err = db.Ping()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database ping failed: %v", err), http.StatusInternalServerError)
		log.Printf("Database ping failed: %v", err)
		return
	}

	// Prepare the query to select all data from the specified table.
	// Using a static string for the table name is safe as it's not from user input.
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	
	// Execute the query.
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	// Get the column names to use as keys for the JSON objects.
	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get columns: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to get columns: %v", err)
		return
	}

	// Create a slice to hold all the Monarch records.
	records := make([]MonarchRecord, 0)
	
	// Prepare a slice to hold the raw values from the database.
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Iterate through the rows and populate the records slice.
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			log.Printf("Failed to scan row: %v", err)
			return
		}

		record := make(MonarchRecord)
		for i, col := range columns {
			// Handle different data types.
			val := values[i]
			if b, ok := val.([]byte); ok {
				// Convert byte slices to strings.
				record[col] = string(b)
			} else {
				record[col] = val
			}
		}
		records = append(records, record)
	}
	
	// Set the content type to application/json.
	w.Header().Set("Content-Type", "application/json")
	
	// Marshal the slice of records into a JSON byte slice.
	jsonData, err := json.Marshal(records)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	// Write the JSON response to the client.
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CHQ: Gemini AI generated function
// helloHandler is the function that will be executed for requests to the "/" route.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "This is the server for the student records app. It's written in Go (aka GoLang).")
}


// faviconHandler serves the favicon.ico file.
func faviconHandler(w http.ResponseWriter, r *http.Request) {
    // Open the favicon file
    favicon, err := os.ReadFile("./static/calculator.ico")
    if err != nil {
        http.NotFound(w, r)
        return
    }

    // Set the Content-Type header
    w.Header().Set("Content-Type", "image/x-icon")
    
    // Write the file content to the response
    w.Write(favicon)
}

func main() {
	// Set up the HTTP router.
		// Initialize the router
	router := mux.NewRouter()

	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/favicon.ico", faviconHandler)

	router.HandleFunc("/api/monarchs", getMonarchsHandler)
	// router.HandleFunc("/api/monarchs", getMonarchsHandler)

	
	// --- CORS Setup ---
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsRouter := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)

	// // Start the server on port 5000.
	// fmt.Println("Server is running on port 5000...")
	// log.Fatal(http.ListenAndServe(":5000", nil))
	
	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	fmt.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}