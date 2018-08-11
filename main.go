package main

import (
	"fmt"
	"log"

	"github.com/mndailey/mws-sales-go/db"
)

func main() {
	dbi, err := db.Instance()
	if err != nil {
		log.Fatal(err)
	}
	if err = dbi.LoadInventory(); err != nil {
		log.Fatal(err)
	}
	defer dbi.Close()
	cnt := dbi.DumpInvMap(func(inv *db.InvRec) string {
		if inv.SkuGrp == "" {
			return fmt.Sprintf("%s -> %s", inv.Sku, inv.SkuGrp)
		} else {
			return fmt.Sprintf("%s -> %s", inv.Sku, inv.SkuGrp)
		}
	})
	fmt.Println("Inv: ", cnt)
}
