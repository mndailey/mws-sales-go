package main

import (
	"github.com/mndailey/mws-sales-go/db"
  "github.com/mndailey/mws-sales-go/util"
)

func main() {
  db.Instance();
  util.WaitForExit()
}
