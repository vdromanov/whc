package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mxk/go-sqlite/sqlite3" //We need <init> func of sql driver
)

type DbRepr struct {
	DbConnStr      string
	TableName      string
	PersonIDColumn string
	IoTimeColumn   string
}

//GetUserIoTimesBetween searches <userId>'s activity between <startUtime> and <endUtime>"
func (sq *DbRepr) GetUserIoTimesBetween(startUtime, endUtime int64, userID int) ([]int64, error) {
	times := make([]int64, 0)

	db, err := sql.Open("sqlite3", sq.DbConnStr)
	if err != nil {
		return times, err
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? AND %s > ? AND %s < ?", sq.IoTimeColumn, sq.TableName, sq.PersonIDColumn, sq.IoTimeColumn, sq.IoTimeColumn)
	rows, err := db.Query(query, userID, startUtime, endUtime)
	if err != nil {
		return times, err
	}
	defer rows.Close()

	for rows.Next() {
		var utime int64
		if err := rows.Scan(&utime); err != nil {
			return times, err
		}
		times = append(times, utime)
	}
	return times, nil
}
