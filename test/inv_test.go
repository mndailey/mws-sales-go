package test

import (
	"fmt"
	"testing"

	"github.com/mndailey/mws-sales-go/db"
)

// TestInv - Test The basic inventory
func TestInv(t *testing.T) {
	dbi, err := db.Instance()
	if err != nil {
		t.Errorf("Error opening db %v.", err)
	} else {
		defer dbi.Close()
		if err = dbi.LoadInventory(); err != nil {
			t.Errorf("Error Loading inventory %v.", err)
		} else {
			cnt := dbi.DumpInvMap(func(inv *db.InvRec) string {
				if inv.TotalQty > 0 {
					return fmt.Sprintf("SKU: %s, GRP: %s, Tot: %d, Stock: %d", inv.Sku, inv.SkuGrp, inv.TotalQty, inv.InStockQty)
				}
				return ""
			})
			fmt.Println("Inv: ", cnt)
		}
	}

}

// TestInv - Test The basic inventory
func TestOrd(t *testing.T) {
	dbi, err := db.Instance()
	if err != nil {
		t.Errorf("Error opening db %v.", err)
	} else {
		defer dbi.Close()
		if err = dbi.LoadOrderTable(); err != nil {
			t.Errorf("Error Loading Order %v.", err)
		} else {
			cnt := dbi.DumpOrderMap(func(sku string, idx, yearweek, qty int) string {
				if qty > 0 {
					return fmt.Sprintf("SKU: %s, Idx: %d, YearWeek: %d, Qty: %d", sku, idx, yearweek, qty)
				}
				return ""
			})
			fmt.Println("Ord: ", cnt)
		}
	}
}
