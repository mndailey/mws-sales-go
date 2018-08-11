package main

import (
	"github.com/mndailey/mws-sales-go/db"
	"log"
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
	for sku, invrec := range dbi.InvMap {
		log.Println(sku, "->", invrec)
	}
}
