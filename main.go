// --- main.go ---

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// The `pq` package is a pure Go PostgreSQL driver for `database/sql`.
	"github.com/descope/go-sdk/descope/client"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// MonarchRecord represents a row from the database table.
// Using a map[string]interface{} is flexible since the schema
// might change, similar to the dynamic dictionary creation in the Python version.
type MonarchRecord map[string]interface{}

// MonarchRecord represents a row from the database table 'june212025'. 
type MyMonarchRecord struct {
	GBIFID                   *string    `json:"gbifID"`
	DatasetKey               *string    `json:"datasetKey"`
	PublishingOrgKey         *string    `json:"publishingOrgKey"`
	EventDate                *string    `json:"eventDate"`
	EventDateParsed          *time.Time `json:"eventDateParsed"`
	Year                     *int       `json:"year"`
	Month                    *int       `json:"month"`
	Day                      *int       `json:"day"`
	DayOfWeek                *int       `json:"day_of_week"`
	WeekOfYear               *int64     `json:"week_of_year"`
	DateOnly                 *string    `json:"date_only"`
	ScientificName           *string    `json:"scientificName"`
	VernacularName           *string    `json:"vernacularName"`
	TaxonKey                 *int64     `json:"taxonKey"`
	Kingdom                  *string    `json:"kingdom"`
	Phylum                   *string    `json:"phylum"`
	Class                    *string    `json:"class"`
	Order                    *string    `json:"order"`
	Family                   *string    `json:"family"`
	Genus                    *string    `json:"genus"`
	Species                  *string    `json:"species"`
	DecimalLatitude          *float64   `json:"decimalLatitude"`
	DecimalLongitude         *float64   `json:"decimalLongitude"`
	CoordinateUncertaintyInMeters *float64 `json:"coordinateUncertaintyInMeters"`
	CountryCode              *string    `json:"countryCode"`
	StateProvince            *string    `json:"stateProvince"`
	IndividualCount          *int64     `json:"individualCount"`
	BasisOfRecord            *string    `json:"basisOfRecord"`
	RecordedBy               *string    `json:"recordedBy"`
	OccurrenceID             *string    `json:"occurrenceID"`
	CollectionCode           *string    `json:"collectionCode"`
	CatalogNumber            *string    `json:"catalogNumber"`
	County                   *string    `json:"county"`
	CityOrTown               *string    `json:"cityOrTown"`
	TimeOnly                 *string    `json:"time_only"` // Storing as string to handle "time without time zone"
}

var db *sql.DB
var descopeClient *client.DescopeClient

var isAnAdmin bool
// Define a custom key type to avoid collisions
type contextKey string

const contextKeyUserID contextKey = "userID"
const contextKeyTeacherID contextKey = "teacherID" // A key for the teacher ID


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
	// records := make([]MonarchRecord, 0)
	records := make([]MyMonarchRecord, 0)

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

		// record := make(MonarchRecord)
		record := make(MyMonarchRecord)
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
	fmt.Fprint(w, "This is the server for the monarch butterflies app. It's written in Go (aka GoLang).")
}


// faviconHandler serves the favicon.ico file.
func faviconHandler(w http.ResponseWriter, r *http.Request) {
    // Open the favicon file
    favicon, err := os.ReadFile("./static/butterfly_net.ico")
    if err != nil {
        http.NotFound(w, r)
        return
    }

    // Set the Content-Type header
    w.Header().Set("Content-Type", "image/x-icon")
    
    // Write the file content to the response
    w.Write(favicon)
}

 // sessionValidationMiddleware is a middleware to validate the Descope session token.
// func sessionValidationMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		sessionToken := r.Header.Get("Authorization")
// 		if sessionToken == "" {
// 			http.Error(w, "Unauthorized: No session token provided", http.StatusUnauthorized)
// 			return
// 		}

// 		sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")

// 		ctx := r.Context()
// 		authorized, token, err := descopeClient.Auth.ValidateSessionWithToken(ctx, sessionToken)
// 		if err != nil || !authorized {
// 			log.Printf("Session validation failed: %v", err)
// 			http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
// 			return
// 		}
// 		if descopeClient.Auth.ValidateRoles(context.Background(), token, []string{"Game Admin"}) {
// 			isAnAdmin = true
// 		} else {
// 			isAnAdmin = false
// 		}

// 		userID := token.ID
// 		// userRole := token.GetTenants()
// 		// userRole := token.GetTenantValue()
// 		// userRole := token.GetTenants()
// 		if userID == "" {
// 			http.Error(w, "Unauthorized: User ID not found in token", http.StatusUnauthorized)
// 			return
// 		}

// 		// For this example, we assume the player ID is the same as the user ID.
// 		// In a real-world app, you would extract this from custom claims in the token.
// 		playerID := userID

// 		// Store the user ID and teacher ID in the request's context
// 		ctxWithUserID := context.WithValue(ctx, contextKeyUserID, userID)
// 		ctxWithIDs := context.WithValue(ctxWithUserID, contextKeyPlayerID, playerID)

// 		next.ServeHTTP(w, r.WithContext(ctxWithIDs))
// 	})
// }
func main() {
	// Set up the HTTP router.
		// Initialize the router
	router := mux.NewRouter()

	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/favicon.ico", faviconHandler)
	// Protected routes (require session validation)
	protectedRoutes := router.PathPrefix("/api").Subrouter()
	// protectedRoutes.Use(sessionValidationMiddleware) // Apply middleware to all routes in this subrouter

	// protectedRoutes.HandleFunc("/monarchsjune2025/{id}", getAllMonarchs).Methods("GET")
	protectedRoutes.HandleFunc("/monarchsjune2025", getAllMonarchs).Methods("GET")


	// router.HandleFunc("/june212025", getMonarchsHandler)
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



func getAllMonarchsAsAdmin(w http.ResponseWriter) {
	var monarchButterflies []MyMonarchRecord
	query := `SELECT * FROM "2025_M06_JUN_2025_butterflies_CT" ORDER BY "date_only"`
	rows, err := db.Query(query)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving butterflies: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var monarchButterfly MyMonarchRecord
		err := rows.Scan(&monarchButterfly.DateOnly, &monarchButterfly.TimeOnly, &monarchButterfly.CityOrTown, &monarchButterfly.County, &monarchButterfly.StateProvince)
		if err != nil {
			log.Printf("Error scanning monarch butterfly row: %v", err)
			continue
		}
		monarchButterflies = append(monarchButterflies, monarchButterfly)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating over monarch butterfly rows: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monarchButterflies)
}

// getAllgodbstudents handles GET requests to retrieve all student records for the authenticated teacher.
// func getAllMonarchsAsTeacher(w http.ResponseWriter, r *http.Request) {
// 	teacherID, ok := r.Context().Value(contextKeyTeacherID).(string)
// 	if !ok || teacherID == "" {
// 		http.Error(w, "Forbidden: Teacher ID not found in session", http.StatusForbidden)
// 		return
// 	}

// 	var students []Student
// 	query := `SELECT id, first_name, last_name, email, major, teacher_id FROM godbstudents WHERE teacher_id = $1 ORDER BY id`
// 	rows, err := db.Query(query, teacherID)

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error retrieving students: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var student Student
// 		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Major, &student.TeacherID)
// 		if err != nil {
// 			log.Printf("Error scanning student row: %v", err)
// 			continue
// 		}
// 		students = append(students, student)
// 	}

// 	if err = rows.Err(); err != nil {
// 		http.Error(w, fmt.Sprintf("Error iterating over student rows: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(students)
// }
func getAllMonarchs(w http.ResponseWriter, r *http.Request){
	// if (isAnAdmin) {
		getAllMonarchsAsAdmin(w)
	// } else {
		// getAllMonarchsAsTeacher(w, r)
	// }
}