package types

import (
	"time"

	"github.com/Canto-Network/Canto/v6/app"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (e *Epoch) Validate() error {
	if e.Duration != types.DefaultUnbondingTime {
		return ErrInvalidEpochDuration
	}
	if !app.EnableAdvanceEpoch {
		if !e.StartTime.Before(time.Now()) {
			return ErrInvalidEpochStartTime
		}
	}
	return nil
}
