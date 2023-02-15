module github.com/Chainflow/oasis-mission-control

replace (
	// Fixes vulnerabilities in etcd v3.3.{10,13} (dependencies via viper).
	// Can be removed once there is a spf13/viper release with updated etcd.
	// https://github.com/spf13/viper/issues/956
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
	// Updates the version used in spf13/cobra (dependency via tendermint) as
	// there is no release yet with the fix. Remove once an updated release of
	// spf13/cobra exists and tendermint is updated to include it.
	// https://github.com/spf13/cobra/issues/1091
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2

	github.com/tendermint/tendermint => github.com/oasisprotocol/tendermint v0.34.9-oasis2
	golang.org/x/crypto/curve25519 => github.com/oasisprotocol/ed25519/extra/x25519 v0.0.0-20210127160119-f7017427c1ea
	golang.org/x/crypto/ed25519 => github.com/oasisprotocol/ed25519 v0.0.0-20210127160119-f7017427c1ea
)

require (
	github.com/fxamacker/cbor v1.5.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/oasisprotocol/oasis-core/go v0.2102.1
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/sendgrid/rest v2.6.4+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.10.0+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.9
	google.golang.org/grpc v1.38.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
)

go 1.15
