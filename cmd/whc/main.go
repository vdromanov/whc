package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/vdromanov/whc/database"
	"github.com/vdromanov/whc/gsheets"
	"github.com/vdromanov/whc/unitime"
)

//Env var names
const (
	credentialsFnameEnv = "CREDENTIALS_FNAME"
	spreadSheetIDEnv    = "SPREADSHEET_ID"
	dbFileEnv           = "DATABASE"
)

var (
	userID           string
	credentialsFname string
	spreadSheetID    string
	dbFname          string
	envVars          = []string{credentialsFnameEnv, spreadSheetIDEnv, dbFileEnv}
)

func getMissingEnvVar(namesSlice []string) string {
	for _, val := range namesSlice {
		if len(os.Getenv(val)) == 0 {
			return val
		}
	}
	return ""
}

func getMinMax(s []int64) (int64, int64) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return s[0], s[len(s)-1]
}

func main() {
	//Set DB table structure
	db := database.DbRepr{
		TableName:      "log_record",
		IoTimeColumn:   "io_time",
		PersonIDColumn: "user_id",
		DbConnStr:      "attend.db",
	}

	flag.StringVar(&userID, "user", "1", "Specify user ID")
	flag.Parse()

	missingEnvVar := getMissingEnvVar(envVars)
	if len(missingEnvVar) != 0 {
		fmt.Printf("env var %s was not set\n", missingEnvVar)
		os.Exit(2)
	}
	db.DbConnStr = os.Getenv(dbFileEnv)
	credentialsFname := os.Getenv(credentialsFnameEnv)
	spreadSheetID := os.Getenv(spreadSheetIDEnv)
	// spreadSheetID := "1CTDX18uqWvFLOsSjzoq5-CJcPl7QkBNaVOt_V_MBG3s"

	//Playing with Google Spreadsheets
	sheet := gsheets.GetSpreadSheet(credentialsFname, spreadSheetID).GetSheetByTitle(userID)
	sheet.AppendRow(0, []string{"Append", "To", "This", "Row"})
	sheet.UpdateRowByCellVal("Append", []string{"Empty", "", "", "Space", "Between"})

	//Playing with time and unixtime conversions
	var UnixTimeExample int64 = 1575990685
	TimeExample := "2019-12-10 10:10:11"
	DefaultFmt := "2006-01-02 15:04:05"
	todayUnix := unitime.UnixTimeToday()
	tomorrowUnix := unitime.UnixTimeTomorrow()
	toUnixExample, _ := unitime.ToUnixTime(DefaultFmt, TimeExample)
	fmt.Printf("From %d to %s\n", UnixTimeExample, unitime.FromUnixTime(DefaultFmt, UnixTimeExample))
	fmt.Printf("From %d to %s\n", UnixTimeExample, unitime.FromUnixTime("15:04", UnixTimeExample))
	fmt.Printf("From %s to %d\n", TimeExample, toUnixExample)
	fmt.Printf("Now unixtime: %d\nTomorrow: %d\n", todayUnix, tomorrowUnix)
	fmt.Printf("Hours diff: %.2f\n", unitime.DeltaHours(todayUnix, tomorrowUnix))

	//Database
	times := db.GetUserIoTimesBetween(todayUnix, tomorrowUnix, 18)
	fmt.Println(times)
}
