/*
 * @Author: @eqto
 * @Date: 2022-05-23 19:31:16
 * @Last Modified by: @eqto
 * @Last Modified time: 2022-06-22 21:08:19
 */

package auth

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/budimanlai/go-wallet-api/internal/app/models"
	"github.com/eqto/api-server"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func JWT(ctx api.Context) error {
	token, ok := ctx.Session().Get(`token`).(*jwt.Token)
	if !ok {
		return ctx.StatusUnauthorized(`JWT token not found`)
	}
	if !token.Valid {
		return ctx.StatusUnauthorized(`invalid JWT token`)
	}
	if payload, ok := token.Claims.(jwt.MapClaims); ok {
		uid, ok := payload[`uid`].(string)
		if !ok {
			return fmt.Errorf("Invalid nasabah id or empty")
		}
		cn, e := ctx.Database()
		if e != nil {
			return errors.Wrap(e, `unable to open database`)
		}

		if len(uid) != 0 {
			var decodedByte, e = base64.StdEncoding.DecodeString(uid)
			if e != nil {
				return ctx.StatusBadRequest(e.Error())
			}
			str := strings.Split(string(decodedByte), ":")
			if len(str) != 2 {
				return ctx.StatusBadRequest("Invalid uid format or empty")
			}

			merchant_id := ctx.Session().GetInt("merchant_id")
			e = models.CheckValidMerchantNasabah(cn, merchant_id, str[0], str[1])
			if e != nil {
				return ctx.StatusBadRequest(e.Error())
			}

			ctx.Session().Put(`nasabah_id`, str[0])
		} else {
			return errors.New("Invalid nasabah id or empty")
		}
	}

	return nil
}

func Session(ctx api.Context) error {
	sess := ctx.Request().Header().Get(`Authorization`)
	if sess == `` {
		return errors.New(`no authorization found`)
	}

	jwtString := strings.Split(sess, "Bearer ")[1]

	token, e := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header[`alg`])
		}
		if payload, ok := t.Claims.(jwt.MapClaims); ok {

			issuer, ok := payload[`iss`].(string)
			if !ok {
				return nil, fmt.Errorf(`invalid iss: %v`, payload[`iss`])
			}
			cn, e := ctx.Database()
			if e != nil {
				return nil, errors.Wrap(e, `unable to open database`)
			}
			p, e := models.GetMerchantByKey(cn, issuer)
			if e != nil {
				return nil, e
			}
			ctx.Session().Put(`merchant_id`, p.ID)

			return []byte(p.AuthKey), nil
		} else {
			return nil, errors.New(`invalid JWT payload`)
		}
	})
	if e != nil {
		return e
	}
	ctx.Session().Put(`token`, token)
	return nil
}
