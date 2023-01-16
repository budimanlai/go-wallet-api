package models

import "github.com/eqto/dbm"

func CheckIfTrxReffIDExists(cn *dbm.Connection, merchant_id int, reff_id string) bool {
	r, e := cn.Get(`SELECT id FROM trx WHERE merchant_id = ? and reff_id = ?`, merchant_id, reff_id)
	if e != nil {
		return false
	}

	if r == nil {
		return false
	} else {
		return true
	}
}
