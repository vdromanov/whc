package gsheets

import (
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var AllowedEmptyRows = 3

func GetSpreadSheet(credentialsFname, spreadsheetID string) spreadsheet.Spreadsheet {
	data, err := ioutil.ReadFile(credentialsFname)
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)
	return sheet
}

func GetSheetByTitle(spreadSheet *spreadsheet.Spreadsheet, title string) *spreadsheet.Sheet {
	sheet, err := spreadSheet.SheetByTitle(title)
	checkError(err)
	return sheet
}

func InsertCell(sheet *spreadsheet.Sheet, row, col int, value string) {
	// Update cell content
	sheet.Update(row, col, value)
	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}

func InsertRow(sheet *spreadsheet.Sheet, startRow, startCol int, values []string) {
	for pos, val := range values {
		InsertCell(sheet, startRow, startCol+pos, val)
	}
}

func getFirstEmptyRow(sheet *spreadsheet.Sheet, allowedEmptyRows int) int {
	empties := 0
	for rowCount, row := range sheet.Rows {
		isEmpty := true
		for _, cell := range row {
			if cell.Value != "" {
				isEmpty = false
				empties = 0
				break
			}
		}
		if isEmpty {
			empties++
			if empties > allowedEmptyRows {
				return rowCount - allowedEmptyRows
			}
		}
	}
	return 0
}

func AppendRow(sheet *spreadsheet.Sheet, rowBeginning int, values []string) {
	rowPos := getFirstEmptyRow(sheet, AllowedEmptyRows)
	InsertRow(sheet, rowPos, rowBeginning, values)
}

func findRowByCellVal(sheet *spreadsheet.Sheet, value string) (rowPos, rowBeginning int) {
	for rowCount, row := range sheet.Rows {
		rowBeginning := 0
		isEmpty := true
		for _, cell := range row {
			if (cell.Value == "") && (isEmpty == true) {
				rowBeginning++
			} else {
				isEmpty = false
			}
			if cell.Value == value {
				return rowCount, rowBeginning
			}
		}
	}
	return 0, 0
}

func UpdateRowByCellVal(sheet *spreadsheet.Sheet, cellValue string, values []string) {
	rowPos, rowBeginning := findRowByCellVal(sheet, cellValue)
	if rowPos != 0 {
		InsertRow(sheet, rowPos, rowBeginning, values)
	} else {
		AppendRow(sheet, rowBeginning, values)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
