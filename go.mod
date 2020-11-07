module github.com/Chainflow/oasis-mission-control

go 1.14

replace (
    github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2

	github.com/tendermint/tendermint => github.com/oasisprotocol/tendermint v0.34.0-rc4-oasis2
	golang.org/x/crypto/curve25519 => github.com/oasisprotocol/ed25519/extra/x25519 v0.0.0-20200819094954-65138ca6ec7c
	golang.org/x/crypto/ed25519 => github.com/oasisprotocol/ed25519 v0.0.0-20200819094954-65138ca6ec7c
)

require (
	github.com/fxamacker/cbor v1.5.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/oasisprotocol/oasis-core/go v0.2012.1
	github.com/sendgrid/rest v2.6.2+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.7.1+incompatible
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.33.8
	google.golang.org/grpc v1.33.2
	gopkg.in/go-playground/validator.v9 v9.31.0
)
