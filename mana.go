package iotago

import (
	"github.com/iotaledger/hive.go/core/safemath"
	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/hive.go/lo"
)

// Mana Structure defines the parameters used in mana calculations.
type ManaStructure struct {
	// BitsCount is the number of bits used to represent Mana.
	BitsCount uint8 `serix:"0,mapKey=bitsCount"`
	// GenerationRate is the amount of potential Mana generated by 1 IOTA in 1 slot.
	GenerationRate uint8 `serix:"1,mapKey=generationRate"`
	// GenerationRateExponent is the scaling of GenerationRate expressed as an exponent of 2.
	GenerationRateExponent uint8 `serix:"2,mapKey=generationRateExponent"`
	// DecayFactors is a lookup table of epoch index diff to mana decay factor (slice index 0 = 1 epoch).
	DecayFactors []uint32 `serix:"3,lengthPrefixType=uint16,mapKey=decayFactors"`
	// DecayFactorsExponent is the scaling of DecayFactors expressed as an exponent of 2.
	DecayFactorsExponent uint8 `serix:"4,mapKey=decayFactorsExponent"`
	// DecayFactorEpochsSum is an integer approximation of the sum of decay over epochs.
	DecayFactorEpochsSum uint32 `serix:"5,mapKey=decayFactorEpochsSum"`
	// DecayFactorEpochsSumExponent is the scaling of DecayFactorEpochsSum expressed as an exponent of 2.
	DecayFactorEpochsSumExponent uint8 `serix:"6,mapKey=decayFactorEpochsSumExponent"`
}

func (m ManaStructure) Equals(other ManaStructure) bool {
	return m.BitsCount == other.BitsCount &&
		m.GenerationRate == other.GenerationRate &&
		m.GenerationRateExponent == other.GenerationRateExponent &&
		lo.Equal(m.DecayFactors, other.DecayFactors) &&
		m.DecayFactorsExponent == other.DecayFactorsExponent &&
		m.DecayFactorEpochsSum == other.DecayFactorEpochsSum &&
		m.DecayFactorEpochsSumExponent == other.DecayFactorEpochsSumExponent
}

type RewardsParameters struct {
	// ValidatorBlocksPerSlot is the number of validation blocks that should be issued by a selected validator per slot during its epoch duties.
	ValidatorBlocksPerSlot uint8 `serix:"0,mapKey=validatorBlocksPerSlot"`
	// ProfitMarginExponent is used for shift operation for calculation of profit margin.
	ProfitMarginExponent uint8 `serix:"1,mapKey=profitMarginExponent"`
	// BootstrappingDuration is the length in epochs of the bootstrapping phase, (approx 3 years).
	BootstrappingDuration EpochIndex `serix:"2,mapKey=bootstrappingDuration"`
	// ManaShareCoefficient is the coefficient used for calculation of initial rewards, relative to the term theta/(1-theta) from the Whitepaper, with theta = 2/3.
	ManaShareCoefficient uint64 `serix:"3,mapKey=manaShareCoefficient"`
	// DecayBalancingConstantExponent is the exponent used for calculation of the initial reward.
	DecayBalancingConstantExponent uint8 `serix:"4,mapKey=decayBalancingConstantExponent"`
	// DecayBalancingConstant needs to be an integer approximation calculated based on chosen DecayBalancingConstantExponent.
	DecayBalancingConstant uint64 `serix:"5,mapKey=decayBalancingConstant"`
	// PoolCoefficientExponent is the exponent used for shifting operation in the pool rewards calculations.
	PoolCoefficientExponent uint8 `serix:"6,mapKey=poolCoefficientExponent"`
}

func (r RewardsParameters) Equals(other RewardsParameters) bool {
	return r.ValidatorBlocksPerSlot == other.ValidatorBlocksPerSlot &&
		r.ProfitMarginExponent == other.ProfitMarginExponent && r.BootstrappingDuration == other.BootstrappingDuration &&
		r.ManaShareCoefficient == other.ManaShareCoefficient &&
		r.DecayBalancingConstantExponent == other.DecayBalancingConstantExponent &&
		r.DecayBalancingConstant == other.DecayBalancingConstant &&
		r.PoolCoefficientExponent == other.PoolCoefficientExponent
}

func (r RewardsParameters) TargetReward(index EpochIndex, api API) (Mana, error) {
	if index > r.BootstrappingDuration {
		return Mana(api.ComputedFinalReward()), nil
	}

	decayedInitialReward, err := api.ManaDecayProvider().RewardsWithDecay(Mana(api.ComputedInitialReward()), index, index)
	if err != nil {
		return 0, ierrors.Errorf("failed to calculate decayed initial reward: %w", err)
	}

	return decayedInitialReward, nil
}

func ManaCost(rmc Mana, workScore WorkScore) (Mana, error) {
	return safemath.SafeMul(rmc, Mana(workScore))
}
