package db

import (
	"fmt"
  "log"
)

// FBARec Structure of the inventory record
type FBARec struct {
	ID        int
	ShipmentName string
	ShipmentID string
  ShipmentStatus string
  TotalUnits int
}

// LoadSkuMap Loads the SKU Map table into memory
func (info *Info) LoadFBASHipments() error {
	if info == nil {
		return fmt.Errorf("info cannot be nil in LoadFBASHipments")
	}
	rows, err := info.db.Query("SELECT id_fba_shipments_tbl, shipment_id" +
      ", shipment_name, shipment_status, total_units" +
      " FROM fba_shipments_tbl WHERE is_processed and shipment_status != 'CLOSED'")
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	m := make(map[string]*FBARec)
	for rows.Next() {
    fbaRec := &FBARec{}
		if err := rows.Scan(&fbaRec.ID, &fbaRec.ShipmentID,
      &fbaRec.ShipmentName, &fbaRec.ShipmentStatus,
      &fbaRec.TotalUnits ); err != nil {
			log.Fatal(err)
		}
		m[fbaRec.ShipmentID] = fbaRec
	}
	info.FBAMap = m
	return nil
}
