package sheets

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	service       *sheets.Service
	spreadsheetID string
}

// NewClient initializes the Google Sheets service client
func NewClient(ctx context.Context, credentialsFile, spreadsheetID string) (*Client, error) {
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}
	return &Client{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

// ReadData reads values from a specific sheet range (e.g. "Sheet1!A1:D10")
func (c *Client) ReadData(ctx context.Context, sheetRange string) ([][]interface{}, error) {
	resp, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, sheetRange).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}
	return resp.Values, nil
}

// WriteData writes/appends values to a specific sheet range
func (c *Client) WriteData(ctx context.Context, sheetRange string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	_, err := c.service.Spreadsheets.Values.Append(c.spreadsheetID, sheetRange, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("unable to write data to sheet: %w", err)
	}
	return nil
}

// UpdateData updates values in a specific sheet range
func (c *Client) UpdateData(ctx context.Context, sheetRange string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	_, err := c.service.Spreadsheets.Values.Update(c.spreadsheetID, sheetRange, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("unable to update data in sheet: %w", err)
	}
	return nil
}

// DeleteRow deletes a specific row index (0-indexed) in a sheet title
func (c *Client) DeleteRow(ctx context.Context, sheetName string, rowIndex int) error {
	// 1. Get spreadsheet metadata to find sheetID by sheetName title
	spreadsheet, err := c.service.Spreadsheets.Get(c.spreadsheetID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve spreadsheet metadata: %w", err)
	}

	var sheetID int64
	found := false
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			sheetID = sheet.Properties.SheetId
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("sheet with title '%s' not found", sheetName)
	}

	// 2. Perform BatchUpdate to delete the dimension row
	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteDimension: &sheets.DeleteDimensionRequest{
					Range: &sheets.DimensionRange{
						SheetId:    sheetID,
						Dimension:  "ROWS",
						StartIndex: int64(rowIndex),
						EndIndex:   int64(rowIndex + 1),
					},
				},
			},
		},
	}

	_, err = c.service.Spreadsheets.BatchUpdate(c.spreadsheetID, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("unable to delete row in sheet: %w", err)
	}

	return nil
}

