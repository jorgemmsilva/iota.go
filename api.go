package iotago

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/crypto/blake2b"

	"github.com/iotaledger/hive.go/serializer/v2/serix"
)

var (
	// ErrMissingProtocolParams is returned when ProtocolParameters are missing for operations which require them.
	ErrMissingProtocolParams = errors.New("missing protocol parameters")

	// internal API instance used to encode/decode objects where protocol parameters don't matter.
	_internalAPI   API
	_internalAPIMu = sync.RWMutex{}
)

func init() {
	_internalAPI = V3API(&ProtocolParameters{})
}

// API handles en/decoding of IOTA protocol objects.
type API interface {
	// Encode encodes the given object to bytes.
	Encode(obj any, opts ...serix.Option) ([]byte, error)
	// Decode decodes the given bytes into object.
	Decode(b []byte, obj any, opts ...serix.Option) (int, error)
	// JSONEncode encodes the given object to its json representation.
	JSONEncode(obj any, opts ...serix.Option) ([]byte, error)
	// JSONDecode decodes the json data into object.
	JSONDecode(jsonData []byte, obj any, opts ...serix.Option) error
	// Underlying returns the underlying serix.API instance.
	Underlying() *serix.API
	// TimeProvider returns the underlying time provider used.
	TimeProvider() *TimeProvider
}

// LatestAPI creates a new API instance conforming to the latest IOTA protocol version.
func LatestAPI(protoParams *ProtocolParameters) API {
	return V3API(protoParams)
}

// calls the internally instantiated API to encode the given object.
//
//nolint:unparam
func internalEncode(obj any, opts ...serix.Option) ([]byte, error) {
	_internalAPIMu.RLock()
	defer _internalAPIMu.RUnlock()
	return _internalAPI.Encode(obj, opts...)
}

// calls the internally instantiated API to decode the given object.
func internalDecode(b []byte, obj any, opts ...serix.Option) (int, error) {
	_internalAPIMu.RLock()
	defer _internalAPIMu.RUnlock()
	return _internalAPI.Decode(b, obj, opts...)
}

// SwapInternalAPI swaps the internally used API of this lib with new.
func SwapInternalAPI(newAPI API) {
	_internalAPIMu.Lock()
	defer _internalAPIMu.Unlock()
	_internalAPI = newAPI
}

// NetworkID defines the ID of the network on which entities operate on.
type NetworkID = uint64

// NetworkIDFromString returns the network ID string's numerical representation.
func NetworkIDFromString(networkIDStr string) NetworkID {
	networkIDBlakeHash := blake2b.Sum256([]byte(networkIDStr))
	return binary.LittleEndian.Uint64(networkIDBlakeHash[:])
}

type protocolAPIContext string

// ProtocolAPIContextKey defines the key to use for a context containing a *ProtocolParameters.
const ProtocolAPIContextKey protocolAPIContext = "protocolParameters"

// ProtocolParameters defines the parameters of the protocol.
type ProtocolParameters struct {
	// The version of the protocol running.
	Version byte `serix:"0,mapKey=version"`
	// The human friendly name of the network.
	NetworkName string `serix:"1,lengthPrefixType=uint8,mapKey=networkName"`
	// The HRP prefix used for Bech32 addresses in the network.
	Bech32HRP NetworkPrefix `serix:"2,lengthPrefixType=uint8,mapKey=bech32Hrp"`
	// The minimum pow score of the network.
	MinPoWScore uint32 `serix:"3,mapKey=minPowScore"`
	// The rent structure used by given node/network.
	RentStructure RentStructure `serix:"4,mapKey=rentStructure"`
	// TokenSupply defines the current token supply on the network.
	TokenSupply uint64 `serix:"5,mapKey=tokenSupply"`
	// GenesisUnixTimestamp defines the genesis timestamp at which the slots start to count.
	GenesisUnixTimestamp uint32 `serix:"6,mapKey=genesisUnixTimestamp"`
	// SlotDurationInSeconds defines the duration of each slot in seconds.
	SlotDurationInSeconds uint8 `serix:"7,mapKey=slotDurationInSeconds"`
	// EpochDurationInSlots defines the amount of slots in an epoch.
	EpochDurationInSlots uint32 `serix:"8,mapKey=epochDurationInSlots"`
	// ManaGenerationRate is the amount of potential Mana generated by 1 IOTA in 1 slot.
	ManaGenerationRate uint8 //`serix:"8,mapKey=manaGenerationRate"`
	// StoredManaDecayFactors is a map of slot index to decay factor.
	StoredManaDecayFactors []float64 //`serix:"9,mapKey=storedManaDecayFactors"`
	// PotentialManaDecayFactors is a map of slot index to decay factor.
	PotentialManaDecayFactors []float64 //`serix:"10,mapKey=potentialManaDecayFactors"`
	// MaxCommitableAge defines the maximum age of a slot to which a block can commit relative to the block timestamp.
	MaxCommitableAge uint32 `serix:"11,mapKey=maxCommitableAge"`
	// StakingUnbondingPeriod defines the unbonding period in epochs before an account can stop staking.
	StakingUnbondingPeriod EpochIndex //`serix:"12,mapKey=stakingUnbondingPeriod"`
}

func (p ProtocolParameters) AsSerixContext() context.Context {
	return context.WithValue(context.Background(), ProtocolAPIContextKey, &p)
}

func (p ProtocolParameters) NetworkID() NetworkID {
	return NetworkIDFromString(p.NetworkName)
}

func (p ProtocolParameters) TimeProvider() *TimeProvider {
	return NewTimeProvider(int64(p.GenesisUnixTimestamp), int64(p.SlotDurationInSeconds), int64(p.EpochDurationInSlots))
}

func (p ProtocolParameters) String() string {
	return fmt.Sprintf("ProtocolParameters: {\n\tVersion: %d\n\tNetwork Name: %s\n\tBech32 HRP Prefix: %s\n\tMinimum PoW Score: %d\n\tRent Structure: %v\n\tToken Supply: %d\n\tGenesis Unix Timestamp: %d\n\tSlot Duration in Seconds: %d\n}",
		p.Version, p.NetworkName, p.Bech32HRP, p.MinPoWScore, p.RentStructure, p.TokenSupply, p.GenesisUnixTimestamp, p.SlotDurationInSeconds)
}

func (p ProtocolParameters) DecayProvider() *DecayProvider {
	return NewDecayProvider(p.ManaGenerationRate, p.StoredManaDecayFactors, p.PotentialManaDecayFactors)
}

// Sizer is an object knowing its own byte size.
type Sizer interface {
	// Size returns the size of the object in terms of bytes.
	Size() int
}
