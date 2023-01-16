package models

import (
	"github.com/eqto/dbm"
	"github.com/pkg/errors"
)

type Nasabah struct {
	ID        string
	Phone     string
	Name      string
	Saldo_idr int
	Saldo_gp  int
}

func CheckValidMerchantNasabah(cn *dbm.Connection, merchant_id int, nasabah_id string, nasabah_token string) error {
	nasabah, e := cn.Get(`SELECT * FROM nasabah_merchant WHERE nasabah_id = ? and merchant_id = ? and nasabah_token = ?`,
		nasabah_id, merchant_id, nasabah_token)
	if e != nil {
		return e
	}

	if nasabah == nil {
		return errors.New("Nasabah not found")
	}

	return nil
}

func GetNasabahByID(cn *dbm.Connection, nasabah_id string) (*Nasabah, error) {
	nasabah, e := cn.Get(`SELECT * FROM nasabah WHERE id = ?`, nasabah_id)
	if e != nil {
		return nil, e
	}

	if nasabah == nil {
		return nil, errors.New("Nasabah not found")
	}

	return &Nasabah{
		ID:        nasabah.String("id"),
		Phone:     nasabah.String("phone"),
		Name:      nasabah.String("fullname"),
		Saldo_idr: nasabah.Int("balance"),
		Saldo_gp:  nasabah.Int("voucher_point"),
	}, nil
}

func IsNasabahExistsByPhone(cn *dbm.Connection, phone string) bool {
	nasabah, e := cn.Get(`SELECT id FROM nasabah WHERE phone = ?`, phone)
	if e != nil {
		return false
	}

	if nasabah == nil {
		return false
	} else {
		return true
	}
}

func IsNasabahExistsByEmail(cn *dbm.Connection, email string) bool {
	nasabah, e := cn.Get(`SELECT id FROM nasabah WHERE email = ?`, email)
	if e != nil {
		return false
	}

	if nasabah == nil {
		return false
	} else {
		return true
	}
}
