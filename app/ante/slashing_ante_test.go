package ante_test

import (
	"github.com/Canto-Network/Canto/v6/app"
	"github.com/Canto-Network/Canto/v6/app/ante"
	"github.com/Canto-Network/Canto/v6/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/evmos/ethermint/encoding"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

var proposer = sdk.AccAddress("test1")

func (suite *AnteTestSuite) TestSlashingParamChangeProposal() {
	suite.SetupTest(false)
	params := suite.app.SlashingKeeper.GetParams(suite.ctx)
	tests := []struct {
		desc                 string
		createSubmitProposal func() *proposal.ParameterChangeProposal
		expectedError        error
	}{
		{
			"SignedBlocksWindow cannot be decreased",
			func() *proposal.ParameterChangeProposal {
				smaller := strconv.FormatInt(params.GetSignedBlocksWindow()-1, 10)
				signedBlocksWindow := proposal.NewParamChange("slashing", "SignedBlocksWindow", smaller)
				return proposal.NewParameterChangeProposal("tc1", "tc1", []proposal.ParamChange{signedBlocksWindow})
			},
			types.ErrInvalidSignedBlocksWindow,
		},
		{
			"SignedBlocksWindow can be increased",
			func() *proposal.ParameterChangeProposal {
				smaller := strconv.FormatInt(params.GetSignedBlocksWindow()+1, 10)
				signedBlocksWindow := proposal.NewParamChange("slashing", "SignedBlocksWindow", smaller)
				return proposal.NewParameterChangeProposal("tc2", "tc2", []proposal.ParamChange{signedBlocksWindow})
			},
			nil,
		},
		{
			"MinSignedPerWindow cannot be decreased",
			func() *proposal.ParameterChangeProposal {
				smaller := params.MinSignedPerWindow.Sub(sdk.OneDec()).String()
				minSignedPerWindow := proposal.NewParamChange("slashing", "MinSignedPerWindow", smaller)
				return proposal.NewParameterChangeProposal("tc3", "tc3", []proposal.ParamChange{minSignedPerWindow})
			},
			types.ErrInvalidMinSignedPerWindow,
		},
		{
			"MinSignedPerWindow can be increased",
			func() *proposal.ParameterChangeProposal {
				smaller := params.MinSignedPerWindow.Add(sdk.OneDec()).String()
				minSignedPerWindow := proposal.NewParamChange("slashing", "MinSignedPerWindow", smaller)
				return proposal.NewParameterChangeProposal("tc4", "tc4", []proposal.ParamChange{minSignedPerWindow})
			},
			nil,
		},
		{
			"DowntimeJailDuration cannot be decreased",
			func() *proposal.ParameterChangeProposal {
				smaller := strconv.FormatInt(int64(params.DowntimeJailDuration)-1, 10)
				downtimeJailDuration := proposal.NewParamChange("slashing", "DowntimeJailDuration", smaller)
				return proposal.NewParameterChangeProposal("tc5", "tc5", []proposal.ParamChange{downtimeJailDuration})
			},
			types.ErrInvalidDowntimeJailDuration,
		},
		{
			"DowntimeJailDuration can be increased",
			func() *proposal.ParameterChangeProposal {
				smaller := strconv.FormatInt(int64(params.DowntimeJailDuration)+1, 10)
				downtimeJailDuration := proposal.NewParamChange("slashing", "DowntimeJailDuration", smaller)
				return proposal.NewParameterChangeProposal("tc6", "tc6", []proposal.ParamChange{downtimeJailDuration})
			},
			nil,
		},
		{
			"SlashFractionDoubleSign cannot be increased",
			func() *proposal.ParameterChangeProposal {
				smaller := params.SlashFractionDoubleSign.Add(sdk.OneDec()).String()
				slashFractionDoubleSign := proposal.NewParamChange("slashing", "SlashFractionDoubleSign", smaller)
				return proposal.NewParameterChangeProposal("tc7", "tc7", []proposal.ParamChange{slashFractionDoubleSign})
			},
			types.ErrInvalidSlashFractionDoubleSign,
		},
		{
			"SlashFractionDoubleSign can be decreased",
			func() *proposal.ParameterChangeProposal {
				smaller := params.SlashFractionDoubleSign.Sub(sdk.OneDec()).String()
				slashFractionDoubleSign := proposal.NewParamChange("slashing", "SlashFractionDoubleSign", smaller)
				return proposal.NewParameterChangeProposal("tc8", "tc8", []proposal.ParamChange{slashFractionDoubleSign})
			},
			nil,
		},
		{
			"SlashFractionDowntime cannot be increased",
			func() *proposal.ParameterChangeProposal {
				smaller := params.SlashFractionDowntime.Add(sdk.OneDec()).String()
				slashFractionDowntime := proposal.NewParamChange("slashing", "SlashFractionDowntime", smaller)
				return proposal.NewParameterChangeProposal("tc9", "tc9", []proposal.ParamChange{slashFractionDowntime})
			},
			types.ErrInvalidSlashFractionDowntime,
		},
	}

	testPrivKeys, _, err := generatePrivKeyAddressPairs(10)
	suite.Require().NoError(err)
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)

	decorator := ante.NewSlashingParamChangeLimitDecorator(&suite.app.SlashingKeeper, suite.app.AppCodec())
	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			msg, err := govtypes.NewMsgSubmitProposal(
				tc.createSubmitProposal(),
				sdk.NewCoins(sdk.NewCoin(suite.app.StakingKeeper.BondDenom(suite.ctx), sdk.NewInt(10000))),
				proposer,
			)
			err = decorator.ValidateMsgs(suite.ctx, []sdk.Msg{msg})
			if tc.expectedError != nil {
				suite.Require().ErrorContains(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
			}

			tx, err := createTx(testPrivKeys[0], []sdk.Msg{msg}...)
			suite.Require().NoError(err)
			txEncoder := encodingConfig.TxConfig.TxEncoder()
			txBytes, err := txEncoder(tx)
			suite.Require().NoError(err)

			resCheckTx := suite.app.DeliverTx(
				abci.RequestDeliverTx{
					Tx: txBytes,
				},
			)
			suite.Require().Equal(resCheckTx.Code, tc.expectedError, resCheckTx.Log)
		})
	}
}
