module github.com/budimanlai/midtrans

go 1.19

replace github.com/budimanlai/gowes_wallet_lib => /Users/budimanlai/Documents/projects/go/semolis-v2/gowes_wallet_lib

require (
	github.com/budimanlai/gowes_wallet_lib v0.0.0-00010101000000-000000000000
	github.com/eqto/api-server v0.12.0
	github.com/eqto/config v0.3.1
	github.com/eqto/dbm v0.13.1
	github.com/eqto/go-logger v0.4.0
	github.com/eqto/service v0.1.6
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/pkg/errors v0.9.1
)

require (
	github.com/eqto/go-json v0.4.1
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/midtrans/midtrans-go v1.3.6
)
