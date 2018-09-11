package db

import (
	"errors"
	"fmt"
	"log"
)

// InvRec Structure of the inventory record
type InvRec struct {
	Sku        string
	SkuGrp     string
	LikeSku    string
	ItemName   string
	InReport   bool
	TotalQty   int
	InStockQty int
	FbaEnabled bool
}

// IsFBAEnabled - returns true if FBAEnabled
func (inv *InvRec) IsFBAEnabled() bool {
	if inv != nil {
		return inv.FbaEnabled
	}
	return false
}

// DumpSkuMap Dumps the SKU Map table
func (info *Info) DumpSkuMap() {
	for key, val := range info.SkuMap {
		fmt.Println(key, "->", val)
	}
}

// LoadSkuMap Loads the SKU Map table into memory
func (info *Info) LoadSkuMap() error {
	if info == nil {
		return errors.New("info cannot be nil in LoadSkuMap")
	}
	rows, err := info.db.Query("SELECT sku, sku_grp FROM sku_map_tbl")
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	m := make(map[string]string)
	for rows.Next() {
		var sku string
		var skuGrp string
		if err := rows.Scan(&sku, &skuGrp); err != nil {
			log.Fatal(err)
		}
		m[sku] = skuGrp
	}
	info.SkuMap = m
	return nil
}

// DumpInvMap Dumps the Inv Map table
func (info *Info) DumpInvMap(filter func(inv *InvRec) string) int {
	cnt := 0
	for _, val := range info.InvMap {
		if str := filter(val); str != "" {
			fmt.Println(str)
			cnt++
		}
	}
	return cnt
}

// LoadInventory loads the inventory file into memory
func (info *Info) LoadInventory() error {
	if info == nil {
		return errors.New("info cannot be nil in LoadSkuMap")
	}
	info.LoadSkuMap()
	rows, err := info.db.Query(`SELECT sku, item_name, fba_total_supply_quantity,
    fba_in_stock_supply_quantity, in_report, like_sku FROM inventory_tbl`)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	m := make(map[string]*InvRec)
	for rows.Next() {
		inv := &InvRec{}
		if err := rows.Scan(&inv.Sku, &inv.ItemName, &inv.TotalQty,
			&inv.InStockQty, &inv.InReport, &inv.LikeSku); err != nil {
			return err
		}
		inv.FbaEnabled = inv.TotalQty > 0
		inv.SkuGrp = info.SkuMap[inv.Sku]
		m[inv.Sku] = inv
	}
	info.InvMap = m
	return nil
}
