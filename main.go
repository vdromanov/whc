package main

import "github.com/vdromanov/whc/gsheets"

func main() {
	CredentialsFname := "client_secret.json"
	SpreadSheetID := "1CTDX18uqWvFLOsSjzoq5-CJcPl7QkBNaVOt_V_MBG3s"
	spreadsheet := gsheets.GetSpreadSheet(CredentialsFname, SpreadSheetID)
	sheet := gsheets.GetSheetByTitle(&spreadsheet, "Sheet1")
	gsheets.InsertCell(sheet, 1, 3, "Hello")
	gsheets.InsertRow(sheet, 2, 1, []string{"First", "Second", "Third"})
	gsheets.AppendRow(sheet, 1, []string{"Append", "To", "This", "Row"})
	gsheets.UpdateRowByCellVal(sheet, "To", []string{"Remove", "From", "This", "Row"})
}
