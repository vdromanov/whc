package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
	dateFmt      = "2 January"
	shortDateFmt = "02.01.2006"
	timeFmt      = "15:04"
	datetimeFmt  = "2006-01-02 15:04"
)

var (
	userID           string
	checkingDate     string
	breakDuration    string
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
	if len(s) == 0 {
		return s, errors.New("No data")
	}
	p := make([]int64, 0)
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
	}

	today := time.Now().Format(shortDateFmt) //default date of sync

	//Parsing args & checking required env vars
	flag.StringVar(&userID, "user", "1", "Specify user ID")
	flag.StringVar(&checkingDate, "date", today, "Specify date of sync")
	flag.StringVar(&breakDuration, "break", "15m", "Dinner duration. Ex: 10m, 1.5h, ...")
	flag.Parse()
	if err := checkMissingEnvVar(envVars); err != nil {
		fmt.Println(err.Error())
		fmt.Printf("Env vars:\n\t%s\nMust be specified\n", strings.Join(envVars, "\n\t"))
		os.Exit(1)
	}

	//Range of working time calculations (24h wide)
	startCheckPeriod, err := unitime.GetBeginningOfDay(shortDateFmt, checkingDate)
	if err != nil {
		fmt.Println("Unable to parse time:", err.Error())
		os.Exit(1)
	}
	endCheckPeriod := unitime.TimeWDelay(startCheckPeriod, "24h")
	startCheckUtime := unitime.ToUnixTime(startCheckPeriod)
	endCheckUtime := unitime.ToUnixTime(endCheckPeriod)

	//Fetching data in defined range from DB
	db.DbConnStr = os.Getenv(dbFileEnv)
	id, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println("User ID should be int")
		os.Exit(1)
	}
	times := db.GetUserIoTimesBetween(startCheckUtime, endCheckUtime, id)
	sortedTimes, err := getSortedUtimes(times)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//Calculating times
	day := unitime.FormatFromUnixTime(dateFmt, startCheckUtime)
	workStart := unitime.TimeFromUnix(sortedTimes[0])
	workStartTime := workStart.Format(timeFmt)
	workEndTime := ""
	workedHours := ""
	if len(sortedTimes) > 1 {
		workEnd := unitime.TimeFromUnix(sortedTimes[len(sortedTimes)-1])
		workEndTime = workEnd.Format(timeFmt)
		workedHours = fmt.Sprintf("%.2f", unitime.DeltaHours(unitime.TimeWDelay(workStart, breakDuration), workEnd))
	}

	//Sending to google
	credentialsFname := os.Getenv(credentialsFnameEnv)
	spreadSheetID := os.Getenv(spreadSheetIDEnv) //Spreadsheet must be shared with service account (<credentialsFname>)
	sheet := gsheets.GetSpreadSheet(credentialsFname, spreadSheetID).GetSheetByTitle(userID)
	rowToGoogle := []string{day, workStartTime, breakDuration, workEndTime, workedHours}
	fmt.Printf("Sending to google: %v\n", rowToGoogle)
	sheet.UpdateRowByCellVal(day, rowToGoogle)
}
