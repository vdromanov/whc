package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mxk/go-sqlite/sqlite3"
)

type DbRepr struct {
	DbConnStr      string
	TableName      string
	PersonIDColumn string
	IoTimeColumn   string
}

func (sq *DbRepr) GetUserIoTimesBetween(startUtime, endUtime int64, userId int) []int64 {
	db, err := sql.Open("sqlite3", sq.DbConnStr)
	if err != nil {
		panic(err.Error)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? AND %s >= ? AND %s <= ?", sq.IoTimeColumn, sq.TableName, sq.PersonIDColumn, sq.IoTimeColumn, sq.IoTimeColumn)
	fmt.Println(query)
	rows, err := db.Query(query, userId, startUtime, endUtime)
	if err != nil {
		panic(err.Error)
	}
	defer rows.Close()

	times := make([]int64, 0)

	for rows.Next() {
		var utime int64
		if err := rows.Scan(&utime); err != nil {
			panic(err.Error)
		}
		fmt.Println(utime)
		times = append(times, utime)
		fmt.Println(times)
	}
	return times
}
