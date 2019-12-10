package main

import (
	"fmt"

	"github.com/vdromanov/whc/gsheets"
	"github.com/vdromanov/whc/unitime"
)

func main() {
	//Playing with Google Spreadsheets
	CredentialsFname := "client_secret.json"
	SpreadSheetID := "1CTDX18uqWvFLOsSjzoq5-CJcPl7QkBNaVOt_V_MBG3s"
	sheet := gsheets.GetSpreadSheet(CredentialsFname, SpreadSheetID).GetSheetByTitle("Sheet1") //User id will be here
	sheet.InsertCell(1, 3, "Hello")
	sheet.InsertRow(2, 1, []string{"First", "Second", "Third"})
	sheet.AppendRow(1, []string{"Append", "To", "This", "Row"})
	sheet.UpdateRowByCellVal("To", []string{"Remove", "From", "This", "Row"})

	//Playing with time and unixtime conversions
	var UnixTimeExample int64 = 1575990685
	TimeExample := "2019-12-10 10:10:11"
	DefaultFmt := "2006-01-02 15:04:05"
	todayUnix := unitime.UnixTimeToday()
	tomorrowUnix := unitime.UnixTimeTomorrow()
	toUnixExample, _ := unitime.ToUnixTime(DefaultFmt, TimeExample)
	fmt.Printf("From %d to %s\n", UnixTimeExample, unitime.FromUnixTime(DefaultFmt, UnixTimeExample))
	fmt.Printf("From %s to %d\n", TimeExample, toUnixExample)
	fmt.Printf("Now unixtime: %d\nTomorrow: %d\n", todayUnix, tomorrowUnix)
	fmt.Printf("Hours diff: %.2f\n", unitime.DeltaHours(todayUnix, tomorrowUnix))
}
