package test

import (
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
		dbi.DumpSkuMap()
		dbi.DumpInvMap()
		dbi.DumpOrderMap()
	}

}
