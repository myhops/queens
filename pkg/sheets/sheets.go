package sheets

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings" // Import strings package

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"
)

func NewSheetReader() {

}

// --- Configuration ---
const (
	// Path to your service account key JSON file
	credentialFile = "path/to/your/service-account-key.json" // <<<--- CHANGE THIS

	// The exact name of the spreadsheet file as it appears in Google Drive
	spreadsheetFileName = "Your Spreadsheet Name Here" // <<<--- CHANGE THIS

	// The exact name of the sheet (tab) inside the spreadsheet
	sheetName = "Sheet1" // <<<--- CHANGE THIS (e.g., "Sales Data")

	// The range you want to read in A1 notation (e.g., "A1:C5", "A:A", "A1")
	readRange = "A1:B5" // <<<--- CHANGE THIS
)

// --- Main Function ---
func run() error {
	ctx := context.Background()

	// 1. Authenticate and get HTTP client
	client, err := getAuthenticatedClient(ctx, credentialFile)
	if err != nil {
		return fmt.Errorf("failed to create authenticated client: %w", err)
	}

	// 2. Find the Spreadsheet ID using the Drive API
	spreadsheetID, err := findSpreadsheetIDByName(ctx, client, spreadsheetFileName)
	if err != nil {
		return fmt.Errorf("Failed to find spreadsheet ID: %w", err)
	}
	// fmt.Printf("Found spreadsheet '%s' with ID: %s\n", spreadsheetFileName, spreadsheetID)

	// 3. Read values from the sheet using the Sheets API - Gets [][]string
	stringValues, err := readSheetValues(ctx, client, spreadsheetID, sheetName, readRange)
	if err != nil {
		return fmt.Errorf("Failed to read sheet values: %w", err)
	}

	// 4. Process/Print the returned string values
	// fmt.Printf("\n--- Data Read as [][]string from range %s!%s ---\n", sheetName, readRange)
	if len(stringValues) == 0 {
		// fmt.Println("[No data returned from the specified range]")
	} else {
		for i, row := range stringValues {
			// Example: Print row number and the content of the row (slice of strings)
			fmt.Printf("Row %d: %v\n", i+1, row)
			// You can access individual cells like: value := row[colIndex]
		}
	}
	fmt.Println("--- End of Data ---")
	return nil
}

// --- Authentication ---
// getAuthenticatedClient creates an authenticated HTTP client using a service account key file.
func getAuthenticatedClient(ctx context.Context, credentialPath string) (*http.Client, error) {
	b, err := os.ReadFile(credentialPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file '%s': %v", credentialPath, err)
	}

	// Define required scopes: Drive (to find the file) and Sheets (to read content)
	// Readonly scopes are sufficient for this reading task.
	config, err := google.JWTConfigFromJSON(b,
		drive.DriveReadonlyScope,          // Scope to search/find files
		gsheets.SpreadsheetsReadonlyScope, // Scope to read sheet values
	)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	return client, nil
}

// --- Find Spreadsheet ID (using Drive API) ---
// findSpreadsheetIDByName searches Google Drive for a spreadsheet with the exact given name
// and returns its ID. It requires the service account to have access to the file.
func findSpreadsheetIDByName(ctx context.Context, client *http.Client, name string) (string, error) {
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	// Build the query to find the specific spreadsheet by name and type, excluding trashed files.
	// Ensure the name doesn't contain single quotes or escape them properly for the query.
	escapedName := strings.ReplaceAll(name, "'", "\\'")
	query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.spreadsheet' and trashed = false", escapedName)

	fileList, err := driveService.Files.List().
		Q(query).
		Fields("files(id, name)"). // Request only the ID and name for efficiency
		PageSize(1).               // We expect only one file with this exact name
		Do()

	if err != nil {
		return "", fmt.Errorf("unable to retrieve files via Drive API: %v", err)
	}

	// Check if any file was found
	if len(fileList.Files) == 0 {
		return "", fmt.Errorf("no spreadsheet file found with name: '%s'. Check name and sharing permissions with service account.", name)
	}

	// Warn if multiple files share the exact same name (uncommon but possible)
	if len(fileList.Files) > 1 {
		log.Printf("Warning: Found %d spreadsheets with the name '%s'. Using the first one found (ID: %s).\n", len(fileList.Files), name, fileList.Files[0].Id)
	}

	// Return the ID of the first file found
	return fileList.Files[0].Id, nil
}

// --- Read Sheet Values (using Sheets API) - Returns [][]string ---
// readSheetValues reads the specified range from the given sheet within the spreadsheet
// and returns the data as a two-dimensional slice of strings.
func readSheetValues(ctx context.Context, client *http.Client, spreadsheetID, sheetName, readRange string) ([][]string, error) {
	sheetsService, err := gsheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		// Return nil slice on error
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	// Construct the full range string for the Sheets API (e.g., 'Sheet Name'!A1:B5)
	// Quoting the sheet name is necessary if it contains spaces or special characters.
	quotedSheetName := sheetName
	if strings.ContainsAny(quotedSheetName, " '") { // Check for space or single quote
		// Escape single quotes within the name itself for Sheets API syntax (' becomes '')
		quotedSheetName = "'" + strings.ReplaceAll(quotedSheetName, "'", "''") + "'"
	}
	fullRange := fmt.Sprintf("%s!%s", quotedSheetName, readRange)

	fmt.Printf("Attempting to read range: %s from spreadsheet ID: %s\n", fullRange, spreadsheetID)

	// Call the Sheets API to get values
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, fullRange).Do()
	if err != nil {
		// Provide more specific error messages for common issues
		if strings.Contains(err.Error(), "Unable to parse range") {
			return nil, fmt.Errorf("sheets API error: Invalid sheet name ('%s') or range ('%s')? Check spelling and syntax. Original error: %v", sheetName, readRange, err)
		}
		if strings.Contains(err.Error(), "Requested entity was not found") {
			// This could be the spreadsheet ID or the sheet name within it
			return nil, fmt.Errorf("sheets API error: Spreadsheet ID '%s' not found, sheet '%s' does not exist, or insufficient permissions. Check sharing. Original error: %v", spreadsheetID, sheetName, err)
		}
		// General error fetching data
		return nil, fmt.Errorf("unable to retrieve data from sheet via Sheets API: %v", err)
	}

	// Handle case where the range is valid but empty
	if len(resp.Values) == 0 {
		fmt.Println("No data found in the specified range.")
		return [][]string{}, nil // Return an empty slice, not an error
	}

	// Convert the API response ( [][]interface{} ) to [][]string
	// Determine the maximum number of columns to handle potentially ragged rows from API
	maxCols := 0
	for _, row := range resp.Values {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// Create the string slice and populate it
	stringValues := make([][]string, len(resp.Values))
	for i, row := range resp.Values {
		stringValues[i] = make([]string, maxCols) // Ensure all inner slices have same length for consistency
		for j := 0; j < maxCols; j++ {
			// Check if the current column index exists in the source row and is not nil
			if j < len(row) && row[j] != nil {
				// Convert the cell value (interface{}) to string using fmt.Sprintf
				stringValues[i][j] = fmt.Sprintf("%v", row[j])
			} else {
				// Assign an empty string for empty cells (nil) or columns beyond the row's actual length
				stringValues[i][j] = ""
			}
		}
	}

	return stringValues, nil // Return the successfully converted data
}
