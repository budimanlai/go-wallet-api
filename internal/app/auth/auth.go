/*
 * @Author: @eqto
 * @Date: 2022-05-23 19:31:16
 * @Last Modified by: @eqto
 * @Last Modified time: 2022-06-22 21:08:19
 */

package auth

import (
	"fmt"

	"github.com/eqto/api-server"
	log "github.com/eqto/go-logger"
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
	customerID := ctx.Session().GetInt(`customer_id`)
	if customerID == 0 {
		return ctx.StatusBadRequest(`no customer_id found`)
	}
	return nil

}

func Session(ctx api.Context) error {
	sess := ctx.Request().Header().Get(`Session`)
	if sess != `` {
		token, e := jwt.Parse(sess, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header[`alg`])
			}
			if payload, ok := t.Claims.(jwt.MapClaims); ok {
				fmt.Println(payload)
				return []byte(""), nil
			} else {
				return nil, errors.New(`invalid JWT payload`)
			}
		})
		if e != nil {
			log.D(e)
		}
		ctx.Session().Put(`token`, token)
	}
	return nil
}
