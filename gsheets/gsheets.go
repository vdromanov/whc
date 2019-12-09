package gsheets

import (
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var AllowedEmptyRows = 3

type SpreadSheet struct {
	spreadsheet.Spreadsheet
}

type Sheet struct {
	*spreadsheet.Sheet
}

func GetSpreadSheet(credentialsFname, spreadsheetID string) SpreadSheet {
	data, err := ioutil.ReadFile(credentialsFname)
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)
	return SpreadSheet{sheet}
}

func (sp SpreadSheet) GetSheetByTitle(title string) *Sheet {
	sh, err := sp.SheetByTitle(title) //A pointer to spreadsheet's sheet
	checkError(err)
	return &Sheet{sh}
}

func (sh Sheet) InsertCell(row, col int, value string) {
	// Update cell content
	sh.Update(row, col, value)
	// Make sure call Synchronize to reflect the changes
	err := sh.Synchronize()
	checkError(err)
}

func (sh Sheet) InsertRow(startRow, startCol int, values []string) {
	for pos, val := range values {
		sh.InsertCell(startRow, startCol+pos, val)
	}
}

func (sh Sheet) getFirstEmptyRow(allowedEmptyRows int) int {
	empties := 0
	for rowCount, row := range sh.Rows {
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

func (sh Sheet) AppendRow(rowBeginning int, values []string) {
	rowPos := sh.getFirstEmptyRow(AllowedEmptyRows)
	sh.InsertRow(rowPos, rowBeginning, values)
}

func (sh Sheet) findRowByCellVal(value string) (rowPos, rowBeginning int) {
	for rowCount, row := range sh.Rows {
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

func (sh Sheet) UpdateRowByCellVal(cellValue string, values []string) {
	rowPos, rowBeginning := sh.findRowByCellVal(cellValue)
	if rowPos != 0 {
		sh.InsertRow(rowPos, rowBeginning, values)
	} else {
		sh.AppendRow(rowBeginning, values)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
