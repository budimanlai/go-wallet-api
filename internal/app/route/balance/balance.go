package balance

import (
	"github.com/budimanlai/go-wallet-api/internal/app/models"
	"github.com/eqto/api-server"
	"github.com/eqto/go-json"
)

func Get(ctx api.Context) error {
	db, e := ctx.Database()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	nasabah_id := ctx.Session().GetString("nasabah_id")
	if len(nasabah_id) == 0 {
		return ctx.StatusBadRequest("Invalid nassabah ID or empty")
	}

	nasabah, e := models.GetNasabahByID(db, nasabah_id)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	return ctx.Write(json.Object{
		`balance_idr`: nasabah.Saldo_idr,
		`balance_gp`:  0,
	})
}
