package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"backend/sheets"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	credentialsFile := os.Getenv("CREDENTIALS_FILE")
	if credentialsFile == "" {
		credentialsFile = "credentials.json"
	}

	if spreadsheetID == "" {
		log.Println("WARNING: SPREADSHEET_ID is not set in environment variables")
	}

	ctx := context.Background()
	var sheetsClient *sheets.Client
	var err error

	// We only initialize Google Sheets if credentials file exists and SPREADSHEET_ID is set
	if _, errFile := os.Stat(credentialsFile); errFile == nil && spreadsheetID != "" {
		sheetsClient, err = sheets.NewClient(ctx, credentialsFile, spreadsheetID)
		if err != nil {
			log.Fatalf("Failed to initialize Google Sheets client: %v", err)
		}
		log.Println("Google Sheets client initialized successfully")
	} else {
		log.Println("Google Sheets client skipped: credentials.json or SPREADSHEET_ID missing")
	}

	// Basic CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Routes
	http.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	}))

	http.HandleFunc("/api/data", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if sheetsClient == nil {
			http.Error(w, "Google Sheets database is not configured (missing credentials.json or SPREADSHEET_ID)", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodGet {
			// Read data
			sheetRange := r.URL.Query().Get("range")
			if sheetRange == "" {
				sheetRange = "Sheet1!A1:Z1000" // Default range (changed to Z1000 for more rows)
			}

			data, err := sheetsClient.ReadData(r.Context(), sheetRange)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}

		if r.Method == http.MethodPost {
			// Write data
			var body struct {
				Range  string          `json:"range"`
				Values [][]interface{} `json:"values"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if body.Range == "" || len(body.Values) == 0 {
				http.Error(w, "range and values are required fields", http.StatusBadRequest)
				return
			}

			err := sheetsClient.WriteData(r.Context(), body.Range, body.Values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Data written successfully"})
			return
		}

		if r.Method == http.MethodPatch {
			// Update specific range/row (using PATCH)
			var body struct {
				Range  string          `json:"range"`
				Values [][]interface{} `json:"values"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if body.Range == "" || len(body.Values) == 0 {
				http.Error(w, "range and values are required fields", http.StatusBadRequest)
				return
			}

			err := sheetsClient.UpdateData(r.Context(), body.Range, body.Values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Data updated successfully"})
			return
		}

		if r.Method == http.MethodDelete {
			// Delete specific row index in a sheet
			sheetName := r.URL.Query().Get("sheetName")
			var body struct {
				RowIndex int `json:"rowIndex"` // 0-indexed row number in Google Sheets
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if sheetName == "" {
				http.Error(w, "sheetName parameter is required", http.StatusBadRequest)
				return
			}

			err := sheetsClient.DeleteRow(r.Context(), sheetName, body.RowIndex)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Row deleted successfully"})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
