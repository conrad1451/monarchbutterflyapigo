// --- main.go ---

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	// The `pq` package is a pure Go PostgreSQL driver for `database/sql`.

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Define the structure for the request body
// type ViewCreationRequest struct {
//     PYear     int    `json:"p_year"`
//     PMonth    string `json:"p_month"`
//     PStartDay int    `json:"p_start_day"`
//     PEndDay   int    `json:"p_end_day"`
//     PState    string `json:"p_state"`
// }

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
// var descopeClient *client.DescopeClient

// var isAnAdmin bool
// // Define a custom key type to avoid collisions
// type contextKey string

// const contextKeyUserID contextKey = "userID"
// const contextKeyTeacherID contextKey = "teacherID" // A key for the teacher ID


// getMonarchsHandler fetches all Monarch butterfly data from the database
// and serves it as a JSON response.
// func getMonarchsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Use the same table name as the Python script for consistency.
// 	tableName := "june212025"
	 
// 	connStr := os.Getenv("DIG_OCEAN_DROPLET_DOCKER_PSQL")
	
// 	// Open a connection to the database.
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
// 		log.Printf("Failed to connect to database: %v", err)
// 		return
// 	}
// 	defer db.Close()

// 	// Ping the database to ensure the connection is live.
// 	err = db.Ping()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Database ping failed: %v", err), http.StatusInternalServerError)
// 		log.Printf("Database ping failed: %v", err)
// 		return
// 	}

// 	// Prepare the query to select all data from the specified table.
// 	// Using a static string for the table name is safe as it's not from user input.
// 	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	
// 	// Execute the query.
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
// 		log.Printf("Query failed: %v", err)
// 		return
// 	}
// 	defer rows.Close()

// 	// Get the column names to use as keys for the JSON objects.
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to get columns: %v", err), http.StatusInternalServerError)
// 		log.Printf("Failed to get columns: %v", err)
// 		return
// 	}

// 	// Create a slice to hold all the Monarch records.
// 	// records := make([]MonarchRecord, 0)
// 	records := make([]MyMonarchRecord, 0)

// 	// Prepare a slice to hold the raw values from the database.
// 	values := make([]interface{}, len(columns))
// 	scanArgs := make([]interface{}, len(columns))
// 	for i := range values {
// 		scanArgs[i] = &values[i]
// 	}

// 	// Iterate through the rows and populate the records slice.
// 	for rows.Next() {
// 		err := rows.Scan(scanArgs...)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
// 			log.Printf("Failed to scan row: %v", err)
// 			return
// 		}

// 		// record := make(MonarchRecord)
// 		// CHQ: GEmini AI debugged
// 		// record := make(MyMonarchRecord)
// 		record := MyMonarchRecord{}
// 		for i, col := range columns {
// 			// Handle different data types.
// 			val := values[i]
// 			if b, ok := val.([]byte); ok {
// 				// Convert byte slices to strings.
// 				record[col] = string(b)
// 			} else {
// 				record[col] = val
// 			}
// 		}
// 		records = append(records, record)
// 	}
	
// 	// Set the content type to application/json.
// 	w.Header().Set("Content-Type", "application/json")
	
// 	// Marshal the slice of records into a JSON byte slice.
// 	jsonData, err := json.Marshal(records)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
// 		log.Printf("Failed to marshal JSON: %v", err)
// 		return
// 	}

// 	// Write the JSON response to the client.
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonData)
// }

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
	// protectedRoutes := router.PathPrefix("/api").Subrouter()
	// protectedRoutes.Use(sessionValidationMiddleware) // Apply middleware to all routes in this subrouter

	// protectedRoutes.HandleFunc("/monarchbutterlies/dayscan/{calendarDate}", getSingleDayScan).Methods("GET")
	// protectedRoutes.HandleFunc("/monarchsjune2025", getAllMonarchs).Methods("GET")

	router.HandleFunc("/monarchbutterlies/dayscan/{calendarDate}", getSingleDayScan).Methods("GET")
	router.HandleFunc("/monarchsjune2025", getAllMonarchs).Methods("GET")


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

func getAllMonarchsAsAdmin2(w http.ResponseWriter, _ *http.Request) {
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

// CHQ: Gemini AI corrected function
// Corrected getAllMonarchsAsAdmin to ignore the 'r' parameter
// func getMonarchButterfliesSingleDayAsAdmin(theTablename string, w http.ResponseWriter, r *http.Request) {

func getMonarchButterfliesSingleDayAsAdmin(theTablename string, w http.ResponseWriter, _ *http.Request) {
	// 1. Establish DB Connection
	connStr := os.Getenv("DIG_OCEAN_DROPLET_DOCKER_PSQL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err) // Log 1: Connection failure
		http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 2. Ping DB
	err = db.Ping()
	if err != nil {
		log.Printf("Database ping failed: %v", err) // Log 2: Ping failure
		http.Error(w, fmt.Sprintf("Database ping failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	var monarchButterflies []MyMonarchRecord
	tableName := theTablename
	
	// Explicitly listing all 35 columns to match the struct fields.
	query := fmt.Sprintf(`SELECT "gbifID", "datasetKey", "publishingOrgKey", "eventDate", "eventDateParsed", "year", "month", "day", "day_of_week", "week_of_year", "date_only", "scientificName", "vernacularName", "taxonKey", "kingdom", "phylum", "class", "order", "family", "genus", "species", "decimalLatitude", "decimalLongitude", "coordinateUncertaintyInMeters", "countryCode", "stateProvince", "individualCount", "basisOfRecord", "recordedBy", "occurrenceID", "collectionCode", "catalogNumber", "county", "cityOrTown", "time_only" FROM "%s" ORDER BY "date_only"`, tableName)
	
	// 3. Execute Query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving butterflies: %v", err), http.StatusInternalServerError)
		log.Printf("Query failed for table %s: %v", tableName, err) // Log 3: Query failure
		return
	}
	defer rows.Close()

	// 4. Iterate and Scan Rows
	for rows.Next() {
		var record MyMonarchRecord
		err := rows.Scan(
			&record.GBIFID,
			&record.DatasetKey,
			&record.PublishingOrgKey,
			&record.EventDate,
			&record.EventDateParsed,
			&record.Year,
			&record.Month,
			&record.Day,
			&record.DayOfWeek,
			&record.WeekOfYear,
			&record.DateOnly,
			&record.ScientificName,
			&record.VernacularName,
			&record.TaxonKey,
			&record.Kingdom,
			&record.Phylum,
			&record.Class,
			&record.Order,
			&record.Family,
			&record.Genus,
			&record.Species,
			&record.DecimalLatitude,
			&record.DecimalLongitude,
			&record.CoordinateUncertaintyInMeters,
			&record.CountryCode,
			&record.StateProvince,
			&record.IndividualCount,
			&record.BasisOfRecord,
			&record.RecordedBy,
			&record.OccurrenceID,
			&record.CollectionCode,
			&record.CatalogNumber,
			&record.County,
			&record.CityOrTown,
			&record.TimeOnly,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			log.Printf("Failed to scan row from table %s: %v", tableName, err) // Log 4: Scan failure
			return
		}
		monarchButterflies = append(monarchButterflies, record)
	}

	// 5. Check for Row Iteration Errors
	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating over monarch butterfly rows: %v", err), http.StatusInternalServerError)
		log.Printf("Error iterating over rows from table %s: %v", tableName, err) // Log 5: Row iteration error
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(monarchButterflies)
}

// Corrected getAllMonarchsAsAdmin to ignore the 'r' parameter
func getAllMonarchsAsAdmin(w http.ResponseWriter, _ *http.Request) {
	connStr := os.Getenv("DIG_OCEAN_DROPLET_DOCKER_PSQL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database ping failed: %v", err), http.StatusInternalServerError)
		log.Printf("Database ping failed: %v", err)
		return
	}

	// getMonarchButterfliesSingleDayAsAdmin("june212025", w, nil)
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

// CHQ: Gemini AI corrected parameters to ignore the r
func getAllMonarchs(w http.ResponseWriter, _ *http.Request) {
	// if (isAnAdmin) {
		getAllMonarchsAsAdmin2(w, nil) // You can pass nil as the request since the function doesn't use it

    // getAllMonarchsAsAdmin(w, nil) // You can pass nil as the request since the function doesn't use it
	// } else {
		// getAllMonarchsAsTeacher(w, r)
	// }
}

func generateTableName(day int, monthInt int, year int) string {
	// 1. Define the equivalent of my_calendar (a map in Go)
	// You would typically define this map globally or pass it in,
	// but defining it here works for a direct conversion.

	monthIntToStr := ""

	if(monthInt < 10){
		monthIntToStr = ("0" + strconv.Itoa(monthInt))
	} else {
		monthIntToStr = strconv.Itoa(monthInt)
	}
  
	myCalendar := map[string]string{
		"01":   "january",
		"02":   "february",
		"03":   "march",
		"04":   "april",
		"05":   "may",
		"06":   "june",
		"07":   "july",
		"08":    "august",
		"09": "september",
		"10":   "october",
		"11":  "november",
		"12":  "december",
	}

	// Retrieve the month string from the map
	// monthStr, ok := myCalendar[month]
	monthStr, ok := myCalendar[monthIntToStr]
 
	if !ok {
		// Handle case where the month is not found (optional, but good practice)
		return "Error: Invalid month"
	}

	// 2. Implement the conditional logic and string concatenation
	var tableName string
	yearStr := strconv.Itoa(year) // Convert int year to string

	tableName = fmt.Sprintf("%s%d%s", monthStr, day, yearStr)

	return tableName
}

// CHQ: Gemini AI added log statements to debug
func getSingleDayScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get the date as a string (MMDDYYYY)
	dateStr := vars["calendarDate"]

	// **********************************************
	// DEBUG LOGGING ADDED HERE TO SEE THE RECEIVED DATE STRING AND LENGTH
	// **********************************************
	log.Printf("Received calendarDate: %s (Length: %d)", dateStr, len(dateStr))

	// 1. String length check (must be exactly 8 characters)
	if len(dateStr) != 8 {
		http.Error(w, "Invalid date given - expected 8 digits in MMDDYYYY format", http.StatusBadRequest)
		log.Printf("Invalid date string length: %s, expected 8", dateStr)
		return
	}

	// 2. Extract components using string slicing (MMDDYYYY)
	monthStr := dateStr[0:2] // MM (e.g., "06")
	dayStr := dateStr[2:4]   // DD (e.g., "30")
	yearStr := dateStr[4:8]  // YYYY (e.g., "2025")

	// 3. Convert day and year to integers for table name generation
	dayInt, err := strconv.Atoi(dayStr)
	if err != nil {
		http.Error(w, "Invalid day format in date", http.StatusBadRequest)
		log.Printf("Invalid day format: %s", dayStr)
		return
	}

	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year format in date", http.StatusBadRequest)
		log.Printf("Invalid year format: %s", yearStr)
		return
	}
	
	// Check if monthStr is a valid two-digit number
	// Although we use monthStr in generateTableName, we need to ensure it's a number
	if _, err := strconv.Atoi(monthStr); err != nil {
		http.Error(w, "Invalid month format in date: not a number", http.StatusBadRequest)
		log.Printf("Invalid month format: %s", monthStr)
		return
	}
	
	// The `useVariable` flag is preserved from your original code
	useVariable := false 

	// Generate the dynamic table name using the string-based month
	myChoice := generateTableName(dayInt, monthStr, yearInt)

	// If useVariable is false, override with the hardcoded test table name
	if !useVariable {
		myChoice = "june212025"
		log.Printf("Using hardcoded table name: %s", myChoice)
	} else {
		log.Printf("Using generated table name: %s", myChoice)
	}

	// Call the function to fetch data from the determined table
	getMonarchButterfliesSingleDayAsAdmin(myChoice, w, nil)
}




// createDailyViewHandler processes the POST request to create the view
// func createDailyViewHandler(w http.ResponseWriter, r *http.Request) {
//     // 1. Only allow POST requests
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//         return
//     }

//     // 2. Decode the JSON request body
//     var req ViewCreationRequest
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid JSON input: "+err.Error(), http.StatusBadRequest)
//         return
//     }

//     // 3. Validate required integer fields (string fields will be empty if not provided)
//     if req.PYear == 0 || req.PStartDay == 0 || req.PEndDay == 0 || req.PMonth == "" || req.PState == "" {
//         http.Error(w, "Missing one or more required parameters (p_year, p_month, p_start_day, p_end_day, p_state)", http.StatusBadRequest)
//         return
//     }

//     // 4. Construct the SQL function call
//     // Note: The '$1', '$2', etc., notation is for PostgreSQL parameterized queries,
//     // which prevents SQL injection.
//     sqlCall := `SELECT create_daily_data_view($1, $2, $3, $4, $5)`

//     // 5. Execute the function
//     var confirmationMessage string
//     err := db.QueryRow(
//         sqlCall, 
//         req.PYear, 
//         req.PMonth, 
//         req.PStartDay, 
//         req.PEndDay, 
//         req.PState,
//     ).Scan(&confirmationMessage)

//     if err != nil {
//         // If the database function raised an exception (e.g., invalid day range),
//         // the error will be caught here.
//         log.Printf("Error executing function: %v", err)
//         http.Error(w, fmt.Sprintf("Database operation failed: %s", err.Error()), http.StatusInternalServerError)
//         return
//     }

//     // 6. Send the success response
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     json.NewEncoder(w).Encode(map[string]string{
//         "status": "success",
//         "message": confirmationMessage,
//     })
// }