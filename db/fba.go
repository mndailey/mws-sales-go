package db

import (
	"fmt"
	"log"
	"strconv"
)

// FBADetailRec Structure of the FBA Shipement Detail record
type FBADetailRec struct {
	ID       int
	FBARecID int
	Sku      string
	Received int
	Shipped  int
}

// FBARec Structure of the FBA Shipement record
type FBARec struct {
	ID             int
	ShipmentName   string
	ShipmentID     string
	ShipmentStatus string
	TotalUnits     int
	Detail         []*FBADetailRec
}

//MariaDB [ey_order_process]> SELECT * FROM  LIMIT 1;
//| id_fba_shipments_detail_tbl | id_fba_shipments_tbl | sku     | received | shipped | version |

// LoadFBASHipmentsDetail Loads the SKU Map table into memory
func (info *Info) LoadFBASHipmentsDetail(clause string) ([]*FBADetailRec, error) {
	if info == nil {
		return nil, fmt.Errorf("info cannot be nil in LoadFBASHipments")
	}
	rows, err := info.db.Query("SELECT id_fba_shipments_detail_tbl" +
		", id_fba_shipments_tbl, sku, received, shipped" +
		" FROM fba_shipments_detail_tbl " + clause)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var a []*FBADetailRec
	for rows.Next() {
		fbaDetailRec := &FBADetailRec{}
		if err := rows.Scan(&fbaDetailRec.ID, &fbaDetailRec.FBARecID,
			&fbaDetailRec.Sku, &fbaDetailRec.Received,
			&fbaDetailRec.Shipped); err != nil {
			log.Fatal(err)
		}
		a = append(a, fbaDetailRec)
	}
	return a, nil
}

// LoadFBASHipments Loads the SKU Map table into memory
func (info *Info) LoadFBASHipments() error {
	if info == nil {
		return fmt.Errorf("info cannot be nil in LoadFBASHipments")
	}
	sql := "SELECT id_fba_shipments_tbl, shipment_id" +
		", shipment_name, shipment_status, total_units" +
		" FROM fba_shipments_tbl WHERE is_processed AND " +
		" (shipment_status != 'CLOSED') AND (shipment_status != 'DELETED')"

	rows, err := info.db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	m := make(map[string]*FBARec)
	m1 := make(map[int]*FBARec)
	for rows.Next() {
		fbaRec := &FBARec{}
		if err := rows.Scan(&fbaRec.ID, &fbaRec.ShipmentID,
			&fbaRec.ShipmentName, &fbaRec.ShipmentStatus,
			&fbaRec.TotalUnits); err != nil {
			log.Fatal(err)
		}
		m[fbaRec.ShipmentID] = fbaRec
		m1[fbaRec.ID] = fbaRec
	}

	keys := ""
	for k := range m1 {
		if keys == "" {
			keys = "WHERE id_fba_shipments_tbl IN ("
		} else {
			keys = keys + ","
		}
		keys = keys + strconv.Itoa(k)
	}
	if keys != "" {
		keys = keys + ")"
	}

	a, err := info.LoadFBASHipmentsDetail(keys)
	if err != nil {
		return err
	}
	for _, d := range a {
		if fbaRec := m1[d.FBARecID]; fbaRec != nil {
			fbaRec.Detail = append(fbaRec.Detail, d)
		}
	}

	info.FBAMap = m
	return nil
}

// DumpOrderMap Dumps the Order Map table with filter
func (info *Info) dumpFBARec(r *FBARec) {
	received := 0
	shipped := 0
	for _, d := range r.Detail {
		received += d.Received
		shipped += d.Shipped
	}
	fmt.Printf("ID: %3d, Shipment: %s, ShipName: %s, Status: %s, Total: %d, Shipped: %d, Received: %d\n",
		r.ID, r.ShipmentID, r.ShipmentName, r.ShipmentStatus, r.TotalUnits, shipped, received)
	for _, d := range r.Detail {
		fmt.Printf("\tSku: %s, Shipped: %d, Received: %d\n", d.Sku, d.Shipped, d.Received)
	}
}

// DumpFBAMap Dumps the FBA Map table with filter
func (info *Info) DumpFBAMap() {
	for _, r := range info.FBAMap {
		info.dumpFBARec(r)
	}
	fmt.Println("FBA: ", len(info.FBAMap))
}
