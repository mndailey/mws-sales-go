package db

import (
	"errors"
	"fmt"
)

// LoadOrderTable Loads Sku Sales into memory
func (info *Info) LoadOrderTable() error {
	if info == nil {
		return errors.New("Cannot setup date on nil db connector")
	}

	if info.InvMap == nil {
		if err := info.LoadInventory(); err != nil {
			return err
		}
	}

	yearweekMap := make(map[int]int)
	info.NumWeeks = 60
	N := 2*info.NumWeeks + 1
	minYearWeek := 99999999
	yearweekIdx := make([]int, N)

	for idx := 0; idx < N; idx++ {
		sql := fmt.Sprintf("SELECT YEARWEEK(DATE_ADD(NOW(), INTERVAL %d WEEK))", idx-info.NumWeeks)
		rows, err := info.db.Query(sql)
		if err != nil {
			return err
		}
		rows.Next()
		var yearweek int
		if err := rows.Scan(&yearweek); err != nil {
			return err
		}
		if yearweek < minYearWeek {
			minYearWeek = yearweek
		}
		yearweekMap[yearweek] = idx
		yearweekIdx[idx] = yearweek

	}
	info.YearWeekToIdxMap = yearweekMap
	info.IdxToYearWeekMap = yearweekIdx
	sql := fmt.Sprintf(`
  SELECT
    sku,
    YEARWEEK(purchase_date),
    SUM(quantity_purchased)
  FROM order_history_tbl
  WHERE fulfillment_channel='AFN' AND YEARWEEK(purchase_date)>=%d
  GROUP BY YEARWEEK(purchase_date), sku
  `, minYearWeek)
	fmt.Printf("SQL: %s\n", sql)
	rows, err := info.db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	om := make(map[string][]int)
	for rows.Next() {
		var sku string
		var yearweek, qty int
		if err := rows.Scan(&sku, &yearweek, &qty); err != nil {
			return err
		}
		if info.InvMap[sku].IsFBAEnabled() {
			m := om[sku]
			if m == nil {
				m = make([]int, N)
				om[sku] = m
			}
			m[yearweekMap[yearweek]] = qty
		}
	}
	info.OrderMap = om
	return nil
}

// DumpOrderMapFilter Dumps the Order Map table with filter
func (info *Info) DumpOrderMapFilter(filter func(sku string, idx, yearweek, qty int) string) int {
	cnt := 0
	for sku, m := range info.OrderMap {
		for idx, qty := range m {
			if str := filter(sku, idx, info.IdxToYearWeekMap[idx], qty); str != "" {
				fmt.Println(str)
				cnt++
			}
		}
	}
	return cnt
}

// DumpOrderMapFilter Dumps the Order Map table with filter
func (info *Info) DumpOrderMap() {
	cnt := info.DumpOrderMapFilter(func(sku string, idx, yearweek, qty int) string {
		if qty > 0 {
			return fmt.Sprintf("SKU: %s, Idx: %d, YearWeek: %d, Qty: %d", sku, idx, yearweek, qty)
		}
		return ""
	})
	fmt.Println("Ord: ", cnt)
}
