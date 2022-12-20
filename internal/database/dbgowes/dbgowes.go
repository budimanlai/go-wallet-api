package dbgowes

import (
	"errors"

	"github.com/eqto/dbm"
)

var (
	stdDB  *dbm.Connection
	stdErr error
)

func Open(driver, host string, port int, username, password, name string) error {
	cn, e := dbm.Connect(dbm.Config{
		DriverName: driver,
		Hostname:   host,
		Port:       port,
		Username:   username,
		Password:   password,
		Name:       name,
	})
	if e != nil {
		stdErr = e
		return e
	}
	stdDB = cn
	stdErr = nil
	return nil
}

func Database() (*dbm.Connection, error) {
	if stdDB == nil {
		if stdErr != nil {
			return nil, stdErr
		}
		return nil, errors.New(`database connection not available`)
	}
	return stdDB, nil
}
