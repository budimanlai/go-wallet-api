package models

import (
	"github.com/eqto/dbm"
	"github.com/pkg/errors"
)

type Merchant struct {
	ID      int
	Name    string
	ApiKey  string
	AuthKey string
}

func GetMerchantByKey(cn *dbm.Connection, key string) (*Merchant, error) {
	merchant, e := cn.Get(`SELECT id, merchant_name, api_key, auth_key FROM merchant WHERE api_key = ? AND status = 'active'`, key)
	if e != nil {
		return nil, e
	}

	if merchant == nil {
		return nil, errors.New("Invalid API Key")
	}

	return &Merchant{
		ID:      merchant.Int(`id`),
		Name:    merchant.String(`merchant_name`),
		ApiKey:  merchant.String(`api_key`),
		AuthKey: merchant.String(`auth_key`),
	}, nil
}
