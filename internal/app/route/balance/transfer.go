package balance

import (
	"time"

	"github.com/budimanlai/go-wallet-api/internal/app/library"
	"github.com/budimanlai/go-wallet-api/internal/app/models"
	"github.com/eqto/api-server"
	"github.com/eqto/go-json"
	"github.com/google/uuid"
)

func Transfer(ctx api.Context) error {
	type transfer_params struct {
		ToNasabah   string `validate:"required"`
		Amount      int    `validate:"required"`
		Description string `validate:"required"`
		TrxID       string `validate:"required"`
	}

	nasabah_id := ctx.Session().GetString("nasabah_id")
	merchant_id := ctx.Session().GetInt("merchant_id")

	db, e := ctx.Database()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	tx, e := ctx.Tx()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	js := ctx.Request().JSON()
	data := &transfer_params{
		ToNasabah:   js.GetString(`to_nasabah`),
		Amount:      js.GetInt(`amount`),
		Description: js.GetString(`description`),
		TrxID:       js.GetString(`trx_id`),
	}

	err := library.Validator().Struct(data)
	if err != nil {
		return ctx.StatusBadRequest(library.FirstError(err))
	}

	// check apakah trx_id sudah ada?
	f := models.CheckIfTrxReffIDExists(db, merchant_id, data.TrxID)
	if f {
		return ctx.StatusBadRequest("Duplicate Trx ID")
	}

	// validate apakah nasabah tujuan ada?
	_, e = models.GetNasabahByID(db, data.ToNasabah)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	// validate apakah saldo nasabah cukup?
	nasabah, e := models.GetNasabahByID(db, nasabah_id)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	if nasabah.Saldo_idr < data.Amount {
		return ctx.StatusBadRequest("Balance not enough")
	}

	id_sender := uuid.New()
	id_recipient := uuid.New()
	internal_trx_id := library.GenerateTrxID()
	sender_trx_id := id_sender.String()
	recipient_trx_id := id_recipient.String()

	// buat transaksi pengirim
	_, e = tx.Insert(`trx`, map[string]interface{}{
		`id`:              sender_trx_id,
		`nasabah_id`:      nasabah_id,
		`merchant_id`:     merchant_id,
		`internal_trx_id`: internal_trx_id,
		`reff_id`:         data.TrxID,
		`mutasi`:          "CR",
		`description`:     data.Description,
		`trx_status`:      "C",
		`trx_amount`:      data.Amount,
		`trx_create_at`:   time.Now(),
		`trx_update_at`:   time.Now(),
	})
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	// kurangi saldo idr pengirim
	_, e = tx.Exec("UPDATE nasabah SET balance = balance - ? WHERE id = ?", data.Amount, nasabah_id)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	// buat transaksi penerima
	_, e = tx.Insert(`trx`, map[string]interface{}{
		`id`:              recipient_trx_id,
		`nasabah_id`:      data.ToNasabah,
		`merchant_id`:     merchant_id,
		`internal_trx_id`: internal_trx_id,
		`reff_id`:         data.TrxID,
		`mutasi`:          "DB",
		`description`:     data.Description,
		`trx_status`:      "C",
		`trx_amount`:      data.Amount,
		`trx_create_at`:   time.Now(),
		`trx_update_at`:   time.Now(),
	})
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	// tambah saldo idr penerima
	_, e = tx.Exec("UPDATE nasabah SET balance = balance + ? WHERE id = ?", data.Amount, data.ToNasabah)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	return ctx.Write(json.Object{
		`to_nasabah`:    data.ToNasabah,
		`amount`:        data.Amount,
		`description`:   data.Description,
		`trx_id`:        data.TrxID,
		`wallet_trx_id`: internal_trx_id,
		`status`:        "success",
	})
}
