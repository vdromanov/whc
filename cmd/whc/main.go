package main

import "github.com/vdromanov/whc/gsheets"

func main() {
	CredentialsFname := "client_secret.json"
	SpreadSheetID := "1CTDX18uqWvFLOsSjzoq5-CJcPl7QkBNaVOt_V_MBG3s"
	sheet := gsheets.GetSpreadSheet(CredentialsFname, SpreadSheetID).GetSheetByTitle("Sheet1") //User id will be here
	sheet.InsertCell(1, 3, "Hello")
	sheet.InsertRow(2, 1, []string{"First", "Second", "Third"})
	sheet.AppendRow(1, []string{"Append", "To", "This", "Row"})
	sheet.UpdateRowByCellVal("To", []string{"Remove", "From", "This", "Row"})
}
