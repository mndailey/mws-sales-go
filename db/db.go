package db

import (
	"database/sql"

	/*
	   mysql driver used to make certin mysql avail to sql.Open
	*/
	_ "github.com/go-sql-driver/mysql"
)

// Info implements a High-level DB instance
type Info struct {
	db               *sql.DB
	SkuMap           map[string]string
	InvMap           map[string]*InvRec
	OrderMap         map[string][]int
	YearWeekToIdxMap map[int]int // Map YearWeek to an index
	IdxToYearWeekMap []int       // Map Index to YearWeek
	NumWeeks         int
}

// Close closes the db connection
func (info *Info) Close() {
	if info == nil {
		return
	}
	info.db.Close()
}

// Instance returns an instance of the db
func Instance() (*Info, error) {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "mystic-sa:acres4@tcp(192.168.0.3)/ey_order_process")
	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}
	// Info implements a High-level DB instance
	info := &Info{
		db: db,
	}
	return info, nil
}
