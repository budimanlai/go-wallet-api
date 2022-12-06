package myapp

import (
	appapi "github.com/budimanlai/go-api-template/internal/app"
	log "github.com/eqto/go-logger"
	"github.com/eqto/service"
)

func main() {
	defer service.HandlePanic()
	service.OnPanic(appapi.Panic.E)
	service.OnStop(appapi.Stop)

	log.SetFile(`runtime/log/` + service.Filename() + `.log`)

	if e := service.Run(appapi.Run); e != nil {
		log.Fatal(e)
	}

}
