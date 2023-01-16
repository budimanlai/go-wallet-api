package account

import (
	"time"

	"github.com/budimanlai/go-wallet-api/internal/app/library"
	"github.com/budimanlai/go-wallet-api/internal/app/models"
	"github.com/eqto/api-server"
	"github.com/eqto/go-json"
	"github.com/google/uuid"
)

func Register(ctx api.Context) error {
	type register_params struct {
		Handphone  string `validate:"required,min=10"`
		Name       string `validate:"required,min=4"`
		Email      string `validate:"required,email"`
		Address    string `validate:"required"`
		City       string `validate:"required"`
		PostalCode string `validate:"required"`
	}

	db, e := ctx.Database()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	tx, e := ctx.Tx()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	js := ctx.Request().JSON()
	data := &register_params{
		Handphone:  js.GetString(`handphone`),
		Name:       js.GetString(`name`),
		Email:      js.GetString(`email`),
		Address:    js.GetString(`address`),
		City:       js.GetString(`city`),
		PostalCode: js.GetString(`postal_code`),
	}

	err := library.Validator().Struct(data)
	if err != nil {
		return ctx.StatusBadRequest(library.FirstError(err))
	}

	// check apakah nomor handphone sudah digunakan?
	found := models.IsNasabahExistsByPhone(db, data.Handphone)
	if found {
		return ctx.StatusBadRequest(`Phone number already registered`)
	}

	// check apakah nomor handphone sudah digunakan?
	found = models.IsNasabahExistsByEmail(db, data.Email)
	if found {
		return ctx.StatusBadRequest(`Email already registered`)
	}

	id := uuid.New()
	merchant_id := ctx.Session().GetInt(`merchant_id`)
	internal_trx_id := library.GenerateTrxID()
	uuid := id.String()
	va_number := "865015" + data.Handphone[len(data.Handphone)-10:]
	wallet_token := library.RandomString(32)

	// save to database
	_, e = tx.Insert(`nasabah`, map[string]interface{}{
		`id`:              uuid,
		`internal_trx_id`: internal_trx_id,
		`va_number`:       va_number,
		`balance`:         0,
		`voucher_point`:   0,
		`fullname`:        data.Name,
		`email`:           data.Email,
		`phone`:           data.Handphone,
		`status`:          "active",
		`create_at`:       time.Now(),
		`create_by`:       nil,
	})

	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	// save token to database
	_, e = tx.Insert(`nasabah_merchant`, map[string]interface{}{
		`nasabah_id`:    uuid,
		`merchant_id`:   merchant_id,
		`nasabah_token`: wallet_token,
		`create_at`:     time.Now(),
	})

	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	return ctx.Write(json.Object{
		`nasabah_id`: uuid,
		`handphone`:  data.Handphone,
		`email`:      data.Email,
		`name`:       data.Email,
		`token`:      wallet_token,
	})
}
