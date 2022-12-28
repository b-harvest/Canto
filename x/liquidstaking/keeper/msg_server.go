package keeper

import "github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}
