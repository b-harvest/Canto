package types_test

import (
	"testing"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type msgTestSuite struct {
	suite.Suite
}

func TestMsgTestSuite(t *testing.T) {
	suite.Run(t, new(msgTestSuite))
}

func (suite *msgTestSuite) TestMsgLiquidStake() {
	delegator := sdk.AccAddress("1")
	stakingCoin := sdk.NewCoin("token", sdk.NewInt(1))

	tcs := []struct {
		desc        string
		expectedErr string
		msg         *types.MsgLiquidStake
	}{
		{
			"happy case",
			"",
			types.NewMsgLiquidStake(delegator.String(), stakingCoin),
		},
		{
			"fail: empty address",
			"invalid delegator address : empty address string is not allowed",
			types.NewMsgLiquidStake("", stakingCoin),
		},
		{
			"fail: zero amount",
			"staking amount must not be zero: invalid request",
			types.NewMsgLiquidStake(delegator.String(), sdk.NewCoin("token", sdk.ZeroInt())),
		},
		{
			"fail: minus amount",
			"negative coin amount: -1",
			types.NewMsgLiquidStake(delegator.String(), sdk.Coin{
				Denom:  "token",
				Amount: sdk.ZeroInt().Sub(sdk.OneInt()),
			}),
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.expectedErr, func() {
			suite.IsType(&types.MsgLiquidStake{}, tc.msg)
			suite.Equal(types.TypeMsgLiquidStake, tc.msg.Type())
			suite.Equal(types.RouterKey, tc.msg.Route())
			suite.Equal(
				sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)),
				tc.msg.GetSignBytes(),
			)

			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				suite.Nil(err)
				signers := tc.msg.GetSigners()
				suite.Len(signers, 1)
				suite.Equal(tc.msg.GetDelegator(), signers[0])
			} else {
				suite.EqualError(err, tc.expectedErr)
			}
		})
	}
}
