package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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

//Formats for time convertions
const (
	dateFmt     = "2 January"
	timeFmt     = "15:04"
	datetimeFmt = "2006-01-02 15:04"
)

var (
	userID           string
	credentialsFname string
	spreadSheetID    string
	dbFname          string
	envVars          = []string{credentialsFnameEnv, spreadSheetIDEnv, dbFileEnv}
)

func checkMissingEnvVar(namesSlice []string) error {
	for _, val := range namesSlice {
		if len(os.Getenv(val)) == 0 {
			return errors.New(fmt.Sprintf("env var %s is missing", val))
		}
	}
	return nil
}

func getSortedUtimes(s []int64) ([]int64, error) {
	p := make([]int64, 0)
	if len(s) == 0 {
		return append(p, 0, 0), errors.New("Nothing to sort")
	}
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	p = append(p, s[0])
	if s[0] < s[len(s)-1] {
		p = append(p, s[len(s)-1])
	}
	return p, nil
}

func main() {
	//Set DB table structure
	db := database.DbRepr{
		TableName:      "log_record",
		IoTimeColumn:   "io_time",
		PersonIDColumn: "user_id",
		DbConnStr:      "attend.db",
	}

	//Parsing args & checking required env vars
	flag.StringVar(&userID, "user", "1", "Specify user ID")
	flag.Parse()
	if err := checkMissingEnvVar(envVars); err != nil {
		fmt.Println(err.Error())
		fmt.Printf("Env vars:\n\t%s\nMust be specified\n", strings.Join(envVars, "\n\t"))
		os.Exit(1)
	}

	db.DbConnStr = os.Getenv(dbFileEnv)
	credentialsFname := os.Getenv(credentialsFnameEnv)
	spreadSheetID := os.Getenv(spreadSheetIDEnv) //Spreadsheet must be shared with service account (<credentialsFname>)

	sheet := gsheets.GetSpreadSheet(credentialsFname, spreadSheetID).GetSheetByTitle(userID)

	todayUnix := unitime.UnixTimeToday()
	tomorrowUnix := unitime.UnixTimeTomorrow()
	id, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println("User ID should be int")
		os.Exit(1)
	}
	times := db.GetUserIoTimesBetween(todayUnix, tomorrowUnix, id)
	sortedTimes, err := getSortedUtimes(times)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//Building a row in gsheet
	today := unitime.FromUnixTime(dateFmt, todayUnix)
	workStartTime := unitime.FromUnixTime(timeFmt, sortedTimes[0])
	workEndTime := ""
	workedHours := ""
	if len(sortedTimes) > 1 {
		workEndTime = unitime.FromUnixTime(timeFmt, sortedTimes[len(sortedTimes)-1])
		workedHours = fmt.Sprintf("%.2f", unitime.DeltaHours(sortedTimes[0], sortedTimes[len(sortedTimes)-1]))
	}
	rowToGoogle := []string{today, workStartTime, "", workEndTime, workedHours}
	fmt.Println(rowToGoogle)
	sheet.UpdateRowByCellVal(today, rowToGoogle)

}
