package testutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/Canto-Network/Canto/v6/testutil/network"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/client/cli"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	suite.cfg = cfg

	// genStateLiquidStaking := types.DefaultGenesisState()
	path, err := os.MkdirTemp("/tmp", "lct-*")
	suite.NoError(err)
	suite.network, err = network.New(suite.T(), path, suite.cfg)
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")
	suite.network.Cleanup()
}

func (suite *IntegrationTestSuite) TestCmdParams() {
	val := suite.network.Validators[0]

	tcs := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"dynamic_fee_rate":{"r0":"0.000000000000000000","u_soft_cap":"0.050000000000000000","u_hard_cap":"0.100000000000000000","u_optimal":"0.090000000000000000","slope1":"0.100000000000000000","slope2":"0.400000000000000000","max_fee_rate":"0.500000000000000000"}}`,
		},
		// TODO: output flag is set to text, but output is still json
		{
			"text output",
			[]string{fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`{"dynamic_fee_rate":{"r0":"0.000000000000000000","u_soft_cap":"0.050000000000000000","u_hard_cap":"0.100000000000000000","u_optimal":"0.090000000000000000","slope1":"0.100000000000000000","slope2":"0.400000000000000000","max_fee_rate":"0.500000000000000000"}}`,
		},
	}
	for _, tc := range tcs {
		tc := tc

		suite.Run(tc.name, func() {
			cmd := cli.CmdQueryParams()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			suite.Require().NoError(err)
			suite.Require().Equal(strings.TrimSpace(tc.expectedOutput), strings.TrimSpace(out.String()))
		})
	}
}

func (suite *IntegrationTestSuite) TestLiquidStaking() {
	vals := suite.network.Validators
	clientCtx := vals[0].ClientCtx
	states := suite.getStates(clientCtx)
	suite.True(states.IsZeroState())
}

func (suite *IntegrationTestSuite) getStates(ctx client.Context) types.NetAmountState {
	var states types.QueryStatesResponse
	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryStates(), []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)})
	suite.NoError(err)
	suite.NoError(suite.cfg.Codec.UnmarshalJSON(out.Bytes(), &states), out.String())
	return states.NetAmountState
}
