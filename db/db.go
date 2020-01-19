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
	FBAMap           map[string]*FBARec
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
	db, _ := sql.Open("mysql", "acres4admin:acres4@tcp(localhost)/ey_order_process")
	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}
	// Info implements a High-level DB instance
	info := &Info{
		db: db,
	}

	if err := info.LoadSkuMap(); err != nil {
		defer db.Close()
		return nil, err
	}

	if err := info.LoadInventory(); err != nil {
		defer db.Close()
		return nil, err
	}

	if err := info.LoadOrderTable(); err != nil {
		defer db.Close()
		return nil, err
	}

	if err := info.LoadFBASHipments(); err != nil {
		defer db.Close()
		return nil, err
	}

	return info, nil
}
