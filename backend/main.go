package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/sheets"
	"strconv"

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
	credentialsJSON := os.Getenv("GOOGLE_CREDENTIALS")
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

	// Initialize using GOOGLE_CREDENTIALS JSON string if set, otherwise fallback to credentials.json file
	if credentialsJSON != "" && spreadsheetID != "" {
		sheetsClient, err = sheets.NewClientFromJSON(ctx, []byte(credentialsJSON), spreadsheetID)
		if err != nil {
			log.Fatalf("Failed to initialize Google Sheets client from GOOGLE_CREDENTIALS JSON env: %v", err)
		}
		log.Println("Google Sheets client initialized successfully using GOOGLE_CREDENTIALS env variable")
		titles, errTitle := sheetsClient.GetSheetTitles(ctx)
		if errTitle == nil {
			log.Printf("[SHEETS DETECTED] List of tabs in spreadsheet: %v", titles)
		} else {
			log.Printf("[SHEETS ERROR] Failed to fetch sheet titles: %v", errTitle)
		}
	} else if _, errFile := os.Stat(credentialsFile); errFile == nil && spreadsheetID != "" {
		sheetsClient, err = sheets.NewClientFromFile(ctx, credentialsFile, spreadsheetID)
		if err != nil {
			log.Fatalf("Failed to initialize Google Sheets client from file: %v", err)
		}
		log.Println("Google Sheets client initialized successfully using credentials file")
		titles, errTitle := sheetsClient.GetSheetTitles(ctx)
		if errTitle == nil {
			log.Printf("[SHEETS DETECTED] List of tabs in spreadsheet: %v", titles)
		} else {
			log.Printf("[SHEETS ERROR] Failed to fetch sheet titles: %v", errTitle)
		}
	} else {
		log.Println("Google Sheets client skipped: credentials (file/env) or SPREADSHEET_ID missing")
	}

	// Basic CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.URL.String())
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

	http.HandleFunc("/api/sheets", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if sheetsClient == nil {
			http.Error(w, "Sheets client not initialized", http.StatusInternalServerError)
			return
		}
		titles, err := sheetsClient.GetSheetTitles(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(titles)
	}))

	type ModuleStat struct {
		Name     string `json:"name"`
		Defect   int    `json:"defect"`
		DemoFail int    `json:"demoFail"`
		Pass     int    `json:"pass"`
		Doing    int    `json:"doing"`
	}

	type DashboardStats struct {
		Global  map[string]int `json:"global"`
		Modules []ModuleStat   `json:"modules"`
	}

	http.HandleFunc("/api/dashboard/stats", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if sheetsClient == nil {
			http.Error(w, "Sheets client not initialized", http.StatusInternalServerError)
			return
		}

		sheetNames := []string{
			"หน้าแดชบอร์ด",
			"หุ้นกู้",
			"พันธบัตร",
			"ไม่อยู่ในความต้องการ",
			"อยู่ในความต้องการ",
			"ทุนเรือนหุ้น",
			"คำร้องขอปรับเปลี่ยนหุ้นรายเดือน",
			"เงินฝาก",
			"ตั๋วสัญญาใช้เงิน",
			"จัดหาเงิน",
			"การเงิน",
			"รายงาน",
			"ตั้งค่า",
		}

		type chanResult struct {
			name string
			rows [][]interface{}
			err  error
		}

		ch := make(chan chanResult, len(sheetNames))

		for _, name := range sheetNames {
			go func(sheetName string) {
				// Query range starts from A4 to read data cells only
				rows, err := sheetsClient.ReadData(r.Context(), sheetName+"!A4:G1000")
				ch <- chanResult{name: sheetName, rows: rows, err: err}
			}(name)
		}

		modules := make([]ModuleStat, 0, len(sheetNames))
		global := map[string]int{
			"pass":     0,
			"defect":   0,
			"demoFail": 0,
			"doing":    0,
		}

		for i := 0; i < len(sheetNames); i++ {
			res := <-ch
			if res.err != nil {
				log.Printf("Error reading sheet %s: %v", res.name, res.err)
				modules = append(modules, ModuleStat{Name: res.name})
				continue
			}

			var passCount, defectCount, demoFailCount, doingCount int
			for _, row := range res.rows {
				if len(row) < 7 || row[0] == nil || fmt.Sprintf("%v", row[0]) == "" {
					continue
				}
				status := "Defect"
				if len(row) > 6 && row[6] != nil {
					status = fmt.Sprintf("%v", row[6])
				}

				switch status {
				case "Pass":
					passCount++
				case "Defect":
					defectCount++
				case "Demo Fail":
					demoFailCount++
				case "Doing", "Ready For Demo", "Demo":
					doingCount++
				default:
					defectCount++
				}
			}

			modules = append(modules, ModuleStat{
				Name:     res.name,
				Defect:   defectCount,
				DemoFail: demoFailCount,
				Pass:     passCount,
				Doing:    doingCount,
			})

			global["pass"] += passCount
			global["defect"] += defectCount
			global["demoFail"] += demoFailCount
			global["doing"] += doingCount
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DashboardStats{
			Global:  global,
			Modules: modules,
		})
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
			var rowIndex int
			hasRowIndex := false

			// 1. Try to read from query parameter first (safer against body stripping)
			if qRowIndex := r.URL.Query().Get("rowIndex"); qRowIndex != "" {
				if val, err := strconv.Atoi(qRowIndex); err == nil {
					rowIndex = val
					hasRowIndex = true
				}
			}

			// 2. Fallback to reading from request body
			if !hasRowIndex {
				var body struct {
					RowIndex int `json:"rowIndex"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
					rowIndex = body.RowIndex
					hasRowIndex = true
				}
			}

			if !hasRowIndex {
				http.Error(w, "rowIndex parameter (query or body) is required and must be a valid integer", http.StatusBadRequest)
				return
			}

			if sheetName == "" {
				http.Error(w, "sheetName parameter is required", http.StatusBadRequest)
				return
			}

			err := sheetsClient.DeleteRow(r.Context(), sheetName, rowIndex)
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
