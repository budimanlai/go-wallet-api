package midtrans

import (
	"fmt"
	"strconv"

	wallet "github.com/budimanlai/gowes_wallet_lib"
	"github.com/budimanlai/midtrans/internal/midtrans_lib"
	"github.com/eqto/api-server"
	"github.com/eqto/go-json"
)

func Callback(ctx api.Context) error {
	tx, e := ctx.Tx()
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	jsonString := ctx.Request().Body()
	jsonObject, e := json.Parse(jsonString)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	orderId := jsonObject.GetStringOr("order_id", "")
	if len(orderId) == 0 {
		return ctx.StatusBadRequest("Invalid paramenter. Required order_id")
	}

	// check apakah order_id ada di database?
	model, e := tx.Get(`SELECT * FROM midtrans WHERE order_id = ?`, orderId)
	if e != nil {
		return ctx.StatusBadRequest(e.Error())
	}

	if model == nil {
		return ctx.StatusBadRequest("Order id not found")
	}

	transactionStatusResp, e := midtrans_lib.Midtrans.CheckTransaction(orderId)
	if e != nil {
		//fmt.Println("Error")
	}

	if transactionStatusResp != nil {
		// 5. Do set transaction status based on response from check transaction status
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionStatusResp.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
		} else if transactionStatusResp.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
		} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
		} else if transactionStatusResp.TransactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
		}

		// check apakah sudah pernah sukses?
		if model.String("fraud_status") == "accept" && model.String("transaction_status") == "settlement" {
			return ctx.Write(json.Object{
				`message`: `OK. Already update`,
			})
		}

		fee_amount := 5000
		gross_amount, e := strconv.ParseFloat(transactionStatusResp.GrossAmount, 64)
		if e != nil {
			fmt.Println("Error:", e.Error())
		}
		net_amount := int(gross_amount) - fee_amount

		raw_json := json.Object{
			`midtrans`: transactionStatusResp,
		}
		// update record
		_, e = tx.Exec(`UPDATE midtrans SET fraud_status = ?, 
			transaction_status = ?, 
			transaction_id = ?,
			transaction_time = ?,
			notif_datetime = now(),
			payment_type = ?,
			gross_amount = ?,
			net_amount = ?,
			fee_amount = ?,
			status_code = ?,
			status_message = ?,
			signature_key = ?,
			raw_json = ?
			WHERE order_id = ?`,
			transactionStatusResp.FraudStatus,
			transactionStatusResp.TransactionStatus,
			transactionStatusResp.TransactionID,
			transactionStatusResp.TransactionTime,
			transactionStatusResp.PaymentType,
			gross_amount,
			net_amount,
			fee_amount,
			transactionStatusResp.StatusCode,
			transactionStatusResp.StatusMessage,
			transactionStatusResp.SignatureKey,
			raw_json.ToString(),
			orderId,
		)
		if e != nil {
			return ctx.StatusBadRequest(e.Error())
		}

		client := json.Object{
			`product_id`:        orderId,
			`order_id`:          orderId,
			`merchant_trx_id`:   model.String(`merchant_trx_id`),
			`net_amount`:        net_amount,
			`fee_amount`:        fee_amount,
			`gross_amount`:      gross_amount,
			`status`:            `success`,
			`callback_response`: transactionStatusResp,
		}

		if transactionStatusResp.FraudStatus == "accept" && transactionStatusResp.TransactionStatus == "settlement" {
			description := "Buy Semolis Package"

			// update balance
			e := wallet.AddFund(tx, model.Int("merchant_id"), model.Int("internal_trx_id"),
				model.String("nasabah_id"), net_amount, description, model.Int("id"))
			if e != nil {
				return ctx.StatusBadRequest(e.Error())
			}
		}

		return ctx.Write(client)
	}
	return ctx.Write(json.Object{
		`message`: `Failed to checkt tranaction status`,
	})
}
