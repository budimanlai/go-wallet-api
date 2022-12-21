package library

import (
	"errors"

	"github.com/eqto/dbm"
)

func GetMerchantURL(tx *dbm.Tx, merchant_id int) (string, error) {
	model, e := tx.Get(`SELECT url_notification FROM merchant WHERE id = ?`, merchant_id)
	if e != nil {
		return "", e
	}

	if model == nil {
		return "", errors.New("Merchant not found")
	}

	return model.String("url_notification"), nil
}
