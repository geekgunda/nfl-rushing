package nflrushing

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const dbName = "thescore"
const dbUser = "root"
const dbPass = "rushingpass"
const rushingStatsTable = "rushingstats"
const maxRecordsPerPage = "1000"

// InitDBConn initializes a connection to MySQL DB
func InitDBConn() (err error) {
	connStr := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		return fmt.Errorf("Failed to connect to db [%s] for user [%s]: %v", dbName, dbUser, err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Failed to ping DB: %v", err)
	}
	return nil
}

// InsertRushingStat inserts a single record in DB
func InsertRushingStat(s Stat) error {
	js, err := json.Marshal(s)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO " + rushingStatsTable
	stmt += "(player,yards,longest,touchdowns,misc) VALUES (?,?,?,?,?)"
	_, err = db.Exec(stmt, s.Player, s.YardsParsed, s.LongestParsed, s.Touchdowns, js)
	return err
}

// FetchFilteredStats takes a filter options Request struct and constructs a SQL query
// that is later used to fetch matching records
func FetchFilteredStats(r Request) (Stats, error) {
	stmt := "SELECT misc from " + rushingStatsTable
	var whereClause, sortClause string
	if r.PlayerFilter != "" {
		whereClause = "player like '%" + r.PlayerFilter + "%'"
	}
	if r.SortFilter != "" && !r.IsDefaultSortSelected {
		if r.IsYardsSortSelected {
			sortClause = "order by yards"
		} else if r.IsLongestSortSelected {
			sortClause = "order by longest"
		} else if r.IsTouchdownsSortSelected {
			sortClause = "order by touchdowns"
		}
	}
	if whereClause != "" {
		stmt += " where " + whereClause
	}
	if sortClause != "" {
		stmt += " " + sortClause
		if r.IsOrderDesc {
			stmt += " desc"
		}
	}
	stmt += " limit " + maxRecordsPerPage
	return fetchStatsForStmt(stmt)
}

// FetchStats pulls all records from DB upto a max limit
// TODO: take maxRecordsPerPage as fn argument and support pagination
func FetchStats() (Stats, error) {
	stmt := "SELECT misc from " + rushingStatsTable + " limit " + maxRecordsPerPage
	return fetchStatsForStmt(stmt)
}

// fetchStatsForStmt takes a SQL statement, queries and parses the resultset into a slice
func fetchStatsForStmt(stmt string) (Stats, error) {
	var ss Stats
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data json.RawMessage
		var stat Stat
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &stat); err != nil {
			return nil, err
		}
		ss = append(ss, stat)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ss, nil
}
