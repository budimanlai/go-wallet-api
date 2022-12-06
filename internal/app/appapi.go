/*
 * @Author: @eqto
 * @Date: 2022-05-23 19:18:02
 * @Last Modified by: @eqto
 * @Last Modified time: 2022-06-22 21:07:04
 */

package appapi

import (
	"errors"

	"github.com/budimanlai/go-api-template/internal/app/auth"
	"github.com/budimanlai/go-api-template/internal/app/route"
	"github.com/budimanlai/go-api-template/internal/logistik/dblogistik"
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
	Panic = log.NewWithFile(`logs/panic.log`)
}

func Run() error {
	port := service.GetInt(`port`)
	if port == 0 {
		return errors.New(`required parameter: --port`)
	}
	if e := config.Open(`configs/app.conf`); e != nil {
		return e
	}

	if e := dblogistik.Open(`mysql`, config.GetOr(`DBLogistik.hostname`, `localhost`), config.GetIntOr(`DBLogistik.port`, 3306),
		config.Get(`DBLogistik.username`), config.Get(`DBLogistik.password`), config.Get(`DBLogistik.name`)); e != nil {
		return e
	}

	svr = api.New()

	if e := svr.OpenDatabase(`mysql`, config.GetOr(`Database.hostname`, `localhost`), config.GetIntOr(`Database.port`, 3306),
		config.Get(`Database.username`), config.Get(`Database.password`), config.Get(`Database.name`)); e != nil {
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
