/*
 * @Author: @eqto
 * @Date: 2022-05-23 19:18:02
 * @Last Modified by: @eqto
 * @Last Modified time: 2022-06-22 21:07:04
 */

package appapi

import (
	"errors"

	"github.com/budimanlai/go-wallet-api/internal/app/auth"
	"github.com/budimanlai/go-wallet-api/internal/app/route"
	"github.com/budimanlai/go-wallet-api/internal/database/dbgowes"
	"github.com/eqto/api-server"
	"github.com/eqto/config"
	_ "github.com/eqto/dbm/driver/mysql"
	log "github.com/eqto/go-logger"
	"github.com/eqto/service"
)

var (
	Panic *log.Logger
	svr   *api.Server
)

func init() {
	Panic = log.NewWithFile(`runtime/logs/panic.log`)
}

func Run() error {
	port := service.GetInt(`port`)
	if port == 0 {
		return errors.New(`required parameter: --port`)
	}
	if e := config.Open(`config/main.conf`); e != nil {
		return e
	}

	if e := dbgowes.Open(`mysql`, config.GetOr(`database.hostname`, `localhost`), config.GetIntOr(`database.port`, 3306),
		config.Get(`database.username`), config.Get(`database.password`), config.Get(`database.database`)); e != nil {
		return e
	}

	svr = api.New()

	if e := svr.OpenDatabase(`mysql`, config.GetOr(`dbwallet.hostname`, `localhost`), config.GetIntOr(`dbwallet.port`, 3306),
		config.Get(`dbwallet.username`), config.Get(`dbwallet.password`), config.Get(`dbwallet.database`)); e != nil {
		return e
	}

	svr.NormalizeFunc(true)
	svr.AddMiddleware(auth.Session)
	svr.AddMiddleware(auth.JWT).Secure()
	route.SetServer(svr)
	return svr.Serve(port)
}

func Stop() {
	if e := svr.Shutdown(); e != nil {
		log.E(e)
	}
}
