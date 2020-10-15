module github.com/Chainflow/oasis-mission-control

go 1.14

replace (
	github.com/tendermint/tendermint => github.com/oasisprotocol/tendermint v0.34.0-rc3-oasis1
	golang.org/x/crypto/curve25519 => github.com/oasisprotocol/ed25519/extra/x25519 v0.0.0-20200528083105-55566edd6df0
	golang.org/x/crypto/ed25519 => github.com/oasisprotocol/ed25519 v0.0.0-20200528083105-55566edd6df0
)

require (
	github.com/fxamacker/cbor v1.5.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/influxdata/influxdb1-client v0.0.0-20200515024757-02f0bf5dbca3
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/oasisprotocol/oasis-core/go v0.0.0-20200811212547-481bbd9cbac0
	github.com/sendgrid/rest v2.6.0+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.6.1+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.8
	github.com/zondax/ledger-oasis-go v0.4.0 // indirect
	google.golang.org/grpc v1.31.0
	gopkg.in/go-playground/validator.v9 v9.31.0
)