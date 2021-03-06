package beacon

import "github.com/oasisprotocol/oasis-core/go/consensus/tendermint/api"

const (
	// AppID is the unique application identifier.
	AppID uint8 = 0x04

	// AppName is the ABCI application name.
	// Run before the scheduler application.
	AppName string = "100_beacon"
)

var (
	// EventType is the ABCI event type for beacon events.
	EventType = api.EventTypeForApp(AppName)

	// QueryApp is a query for filtering events processed by the
	// beacon application.
	QueryApp = api.QueryForApp(AppName)

	// KeyGenerated is the ABCI event attribute key for the new
	// beacons (value is a CBOR serialized beacon.GenerateEvent).
	KeyGenerated = []byte("generated")
)
