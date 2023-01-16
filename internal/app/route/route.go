/*
 * @Author: @eqto
 * @Date: 2022-05-24 08:57:07
 * @Last Modified by: @eqto
 * @Last Modified time: 2022-06-22 21:17:31
 */

package route

import (
	"github.com/budimanlai/go-wallet-api/internal/app/route/account"
	"github.com/budimanlai/go-wallet-api/internal/app/route/balance"
	"github.com/eqto/api-server"
)

func SetServer(svr *api.Server) {
	svr.PostAction(account.Register)

	svr.PostAction(balance.Get).Secure()
	svr.PostAction(balance.Transfer).Secure()
}
