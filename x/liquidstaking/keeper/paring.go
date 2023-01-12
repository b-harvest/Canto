package keeper

import (
	"sort"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: brute-force implementation, optimize
// processing insurance unbonding (it is processed only when there is adequate chunk unbond queue)
// search unbonding insurance in paired status and unpair if found
func resolveUnbondingInsuranceInAlivePairs(k *Keeper, ctx *sdk.Context, currentState *types.State) (newState *types.State) {
	newState = &types.State{}
	newState.ChunkBondRequests = append(types.ChunkBondRequests{}, currentState.ChunkBondRequests...)
	newState.ChunkUnbondRequests = append(types.ChunkUnbondRequests{}, currentState.ChunkUnbondRequests...)
	newState.InsuranceBids = append(types.InsuranceBids{}, currentState.InsuranceBids...)
	newState.InsuranceUnbondRequests = append(types.InsuranceUnbondRequests{}, currentState.InsuranceUnbondRequests...)

	newState.AliveChunks = types.FilterSlice(currentState.AliveChunks, func(aliveChunk types.AliveChunk) bool {
		isUnpaired := false
		newState.InsuranceUnbondRequests = types.FilterSlice(newState.InsuranceUnbondRequests,
			func(req types.InsuranceUnbondRequest) bool {
				if req.AliveChunkId == aliveChunk.Id {
					newState.InsuranceUnbonded = append(newState.InsuranceUnbonded,
						types.InsuranceUnbondRequestedAliveChunk{
							AliveChunk:               aliveChunk,
							InsuranceProviderAddress: req.InsuranceProviderAddress,
						},
					)
					isUnpaired = true
					return false
				}
				return true
			})
		return !isUnpaired
	})

	// can unbonding insurance remain?
	if len(newState.InsuranceUnbondRequests) != 0 {
		k.Logger(*ctx).Info("insurance unbonding request must be cleared")
		panic("insurance unbonding request must be cleared")
	}
	return
}

func resolveChunkBondRequest(k *Keeper,
	ctx *sdk.Context,
	liquidStakingState types.LiquidStakingState,
	chunkBondRequest types.ChunkBondRequest,
) error {
	liquidBondDenom := k.LiquidBondDenom(*ctx)
	liquidStaker, err := sdk.AccAddressFromBech32(chunkBondRequest.Address)
	if err != nil {
		return err
	}
	mintAmount, err := types.NativeTokenToLiquidToken(liquidStakingState, chunkBondRequest.TokenAmount)
	if err != nil {
		return err
	}

	// TODO: need spec, below condition can happen if mint rate has changed since chunk bond request
	// chunkSize := k.GetParams(*ctx).ChunkSize
	// if mintAmount.LT(chunkSize) {
	// 	return types.ErrInvalidTokenAmount
	// }
	mintCoins := sdk.NewCoins(sdk.NewCoin(liquidBondDenom, mintAmount))
	if err := k.bk.MintCoins(*ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(*ctx, types.ModuleName, liquidStaker, mintCoins); err != nil {
		return err
	}

	return nil
}

func resolveChunkUnbondRequest(k *Keeper,
	ctx *sdk.Context,
	liquidStakingState types.LiquidStakingState,
	chunkUnbondRequest *types.ChunkUnbondRequest,
	chunkBondRequests types.ChunkBondRequests,
) error {
	liquidUnstaker, err := sdk.AccAddressFromBech32(chunkUnbondRequest.Address)
	if err != nil {
		return err
	}

	// TODO: this may not secure when chunk size is changed during alive period
	chunkSize := k.GetParams(*ctx).ChunkSize
	unbondTokenAmount := sdk.ZeroInt()
	burnTokenAmount := sdk.ZeroInt()

	numChunkBondRequest := uint64(len(chunkBondRequests))
	if chunkUnbondRequest.NumChunkUnbond < numChunkBondRequest {
		// TODO: correct type
		return types.ErrInvalidChunkUnbondRequestId
	}
	for _, chunkBondRequest := range chunkBondRequests {
		unbondTokenAmount = unbondTokenAmount.Add(chunkBondRequest.TokenAmount)
		burnTokenAmount = burnTokenAmount.Add(chunkSize) // Mul?
	}
	chunkUnbondRequest.NumChunkUnbond -= numChunkBondRequest

	liquidBondDenom := k.LiquidBondDenom(*ctx)
	burnToken := sdk.NewCoin(liquidBondDenom, burnTokenAmount)
	if err := k.bk.BurnCoins(*ctx, types.ModuleName, sdk.NewCoins(burnToken)); err != nil {
		return err
	}

	bondDenom := k.stk.BondDenom(*ctx)
	unbondToken := sdk.NewCoin(bondDenom, unbondTokenAmount)
	if err := k.bk.SendCoinsFromModuleToAccount(*ctx,
		types.ModuleName,
		liquidUnstaker,
		sdk.NewCoins(unbondToken),
	); err != nil {
		return err
	}

	return nil
}

// pairs unbond queue with bond queue. they will be optimized out
func resolveUnbondingChunksAndBondingChunks(k *Keeper, ctx *sdk.Context, state *types.State) *types.State {
	indexUnbondRequest := 0
	indexBondRequest := uint64(0)
	numChunkBondRequest := uint64(len(state.ChunkBondRequests))

	for (indexUnbondRequest < len(state.ChunkUnbondRequests)) && (indexBondRequest < numChunkBondRequest) {
		chunkUnbondRequest := &state.ChunkUnbondRequests[indexUnbondRequest]
		newIndexBondRequest := indexBondRequest + chunkUnbondRequest.NumChunkUnbond
		if newIndexBondRequest > numChunkBondRequest {
			newIndexBondRequest = numChunkBondRequest
		}
		chunkBondRequests := state.ChunkBondRequests[indexBondRequest:newIndexBondRequest]
		// TODO: create current liquid staking state
		liquidStakingState := types.LiquidStakingState{}

		// TODO: FIFO rule?
		if err := resolveChunkUnbondRequest(k, ctx, liquidStakingState, chunkUnbondRequest, chunkBondRequests); err != nil {
			panic(err)
		}
		if chunkUnbondRequest.NumChunkUnbond == 0 {
			k.DeleteChunkUnbondRequest(*ctx, chunkUnbondRequest.Id)
			indexUnbondRequest++
		}
		for _, chunkBondRequest := range chunkBondRequests {
			if err := resolveChunkBondRequest(k, ctx, liquidStakingState, chunkBondRequest); err != nil {
				panic(err)
			}
			k.DeleteChunkBondRequest(*ctx, chunkBondRequest.Id)
		}
		indexBondRequest = newIndexBondRequest
	}
	if indexUnbondRequest < len(state.ChunkUnbondRequests) {
		// update last changed unbond request
		k.SetChunkUnbondRequest(*ctx, state.ChunkUnbondRequests[indexUnbondRequest])
	}
	state.ChunkBondRequests = state.ChunkBondRequests[indexBondRequest:]
	state.ChunkUnbondRequests = state.ChunkUnbondRequests[indexUnbondRequest:]
	return state
}

func rankRankable(ctx *sdk.Context, k *Keeper, rankable []types.Rankable) {
	sort.Slice(rankable, func(i, j int) bool {
		lhs := rankable[i]
		rhs := rankable[j]
		lhsFee, err := calcTotalFeeRate(k, ctx, lhs)
		if err != nil {
			panic("fail to calc fee rate")
		}
		rhsFee, err := calcTotalFeeRate(k, ctx, rhs)
		if err != nil {
			panic("fail to calc fee rate")
		}

		// TODO: need spec to secure determinism
		// temporary strategy: 1. compare types: aliveChunk has higher priority to InsuranceBid (for less redelegation)
		//                     2. compare id   : the lower id has higher priority (FIFO)
		if lhsFee.Equal(rhsFee) {
			switch lhsV := lhs.(type) {
			case *types.AliveChunk:
				switch rhsV := rhs.(type) {
				case *types.AliveChunk:
					return lhsV.Id < rhsV.Id
				case *types.InsuranceBid:
					return true
				default:
					panic("invalid type")
				}
			case *types.InsuranceBid:
				switch rhsV := rhs.(type) {
				case *types.AliveChunk:
					return false
				case *types.InsuranceBid:
					return lhsV.Id < rhsV.Id
				default:
					panic("invalid type")
				}
			default:
				panic("invalid type")
			}
		}
		return lhsFee.LT(rhsFee)
	})
}

func (k *Keeper) rankAliveChunks(ctx *sdk.Context, state *types.State) (ret []types.Rankable) {
	for i, _ := range state.AliveChunks {
		ret = append(ret, &state.AliveChunks[i])
	}
	rankRankable(ctx, k, ret)
	return
}

func (k *Keeper) rankAliveChunksAndInsuranceBid(ctx *sdk.Context, state *types.State) (ret []types.Rankable) {
	for i, _ := range state.AliveChunks {
		ret = append(ret, &state.AliveChunks[i])
	}
	for i, _ := range state.InsuranceBids {
		ret = append(ret, &state.InsuranceBids[i])
	}
	rankRankable(ctx, k, ret)
	return
}

// unpair unbonding chunks
func unpairUnbondingChunksInAliveChunks(k *Keeper, ctx *sdk.Context, state *types.State) *types.State {
	// TODO: optimize unpair candidate selection
	// sort alive pairs
	ranked := k.rankAliveChunks(ctx, state)

	indexAliveChunk := len(ranked) - 1
	indexChunkUnbondReq := 0
	for indexChunkUnbondReq < len(state.ChunkUnbondRequests) && indexAliveChunk >= 0 {
		chunkUnbondRequest := &state.ChunkUnbondRequests[indexChunkUnbondReq]
		aliveChunk, ok := ranked[indexAliveChunk].(*types.AliveChunk)
		if !ok {
			panic("unexpected scenario")
		}

		state.ChunkUnbonded = append(state.ChunkUnbonded, types.ChunkUnbondRequestedAliveChunk{
			AliveChunk: *aliveChunk,
			Address:    chunkUnbondRequest.Address,
		})

		indexAliveChunk--
		chunkUnbondRequest.NumChunkUnbond--
		if chunkUnbondRequest.NumChunkUnbond == 0 {
			k.DeleteChunkUnbondRequest(*ctx, chunkUnbondRequest.Id)
			indexChunkUnbondReq++
		}
	}

	state.ChunkUnbondRequests = state.ChunkUnbondRequests[indexChunkUnbondReq:]

	// TODO: redundant reflection. optimize
	state.AliveChunks = types.AliveChunks{}
	for _, elem := range ranked[:indexAliveChunk+1] {
		state.AliveChunks = append(state.AliveChunks, *elem.(*types.AliveChunk))
	}
	return state
}

func unpairUnbondingChunksInInsuranceUnbonded(k *Keeper, ctx *sdk.Context, state *types.State) *types.State {
	if len(state.InsuranceUnbonded) < len(state.ChunkUnbondRequests) {
		panic("unexpected scenario")
	}
	// TODO: remove duplicated
	indexAliveChunk := 0
	indexChunkUnbondReq := 0
	for indexChunkUnbondReq < len(state.ChunkUnbondRequests) && indexAliveChunk < len(state.InsuranceUnbonded) {
		chunkUnbondRequest := &state.ChunkUnbondRequests[indexChunkUnbondReq]
		insuranceUnbonded := state.InsuranceUnbonded[indexAliveChunk]

		state.ChunkUnbonded = append(state.ChunkUnbonded, types.ChunkUnbondRequestedAliveChunk{
			AliveChunk: insuranceUnbonded.AliveChunk,
			Address:    chunkUnbondRequest.Address,
		})

		indexAliveChunk++
		chunkUnbondRequest.NumChunkUnbond--
		if chunkUnbondRequest.NumChunkUnbond == 0 {
			k.DeleteChunkUnbondRequest(*ctx, chunkUnbondRequest.Id)
			indexChunkUnbondReq++
		}
	}

	state.ChunkUnbondRequests = state.ChunkUnbondRequests[indexChunkUnbondReq:]
	state.InsuranceUnbonded = state.InsuranceUnbonded[indexAliveChunk:]

	return state
}

func (k *Keeper) ResolveUnbondingQueues(ctx sdk.Context, currentState types.State) types.State {
	newState := resolveUnbondingInsuranceInAlivePairs(k, &ctx, &currentState)
	newState = resolveUnbondingChunksAndBondingChunks(k, &ctx, newState)
	newState = unpairUnbondingChunksInAliveChunks(k, &ctx, newState)
	newState = unpairUnbondingChunksInInsuranceUnbonded(k, &ctx, newState)

	// All UnbondRequest must be handled
	if len(newState.InsuranceUnbondRequests) > 0 {
		panic("unexpected scenario")
	}
	if len(newState.ChunkUnbondRequests) > 0 {
		panic("unexpected scenario")
	}
	return *newState
}

func calcTotalFeeRate(k *Keeper, ctx *sdk.Context, elem types.Rankable) (sdk.Dec, error) {
	addr, err := sdk.ValAddressFromBech32(elem.GetValidatorAddress())

	if err != nil {
		return sdk.ZeroDec(), err
	}
	return elem.GetInsuranceFeeRate().Add(k.stk.Validator(*ctx, addr).GetCommission()), nil
}

func (k *Keeper) unpairUnrankedAliveChunks(ctx *sdk.Context, unranked types.Rankables) (ret types.AliveChunks) {
	for _, elem := range unranked {
		aliveChunk, ok := elem.(*types.AliveChunk)
		if !ok {
			continue
		}
		// chunk is unpaired to be paired with ranked insurance bid
		// reuse insurance bid from unranked aliveChunk
		id := k.GetLastInsuranceBidId(*ctx) + 1
		ret = append(ret, *aliveChunk)
		bid := types.InsuranceBid{
			Id:                       id,
			ValidatorAddress:         aliveChunk.ValidatorAddress,
			InsuranceProviderAddress: aliveChunk.InsuranceProviderAddress,
			InsuranceAmount:          aliveChunk.InsuranceAmount,
			InsuranceFeeRate:         aliveChunk.InsuranceFeeRate,
		}
		k.SetLastInsuranceBidId(*ctx, id)
		k.SetInsuranceBid(*ctx, bid)
	}
	return
}

func NewDelegationState(rankedItems types.Rankables, state *types.State) (ret types.DelegationState) {
	ret.DelegationMap = make(map[string]sdk.Int)

	for _, item := range rankedItems {
		found := false
		addr := item.GetValidatorAddress()
		for _, val := range ret.SortedValidators {
			if val.OperatorAddress == addr {
				found = true
				break
			}

		}
		if !found {
			ret.SortedValidators = append(ret.SortedValidators, types.LiquidValidator{
				OperatorAddress: addr,
			})
			ret.DelegationMap[addr] = sdk.ZeroInt()
		}
	}

	for _, item := range state.InsuranceUnbonded {
		found := false
		addr := item.ValidatorAddress
		if !found {
			ret.SortedValidators = append(ret.SortedValidators, types.LiquidValidator{
				OperatorAddress: addr,
			})
			ret.DelegationMap[addr] = sdk.ZeroInt()
		}
	}
	return
}

func (k *Keeper) unpairPairRankedAlivedChunks(ctx *sdk.Context,
	numAliveChunks int64,
	rankedItems types.Rankables,
	state *types.State,
) (types.DelegationState, error) {
	delegationState := NewDelegationState(rankedItems, state)
	unpairedAliveChunks := types.AliveChunks{}
	if len(rankedItems) > int(numAliveChunks) {
		unpairedAliveChunks = k.unpairUnrankedAliveChunks(ctx, rankedItems[numAliveChunks:])
		rankedItems = rankedItems[:numAliveChunks]
	}

	lastAliveChunkId := k.GetLastAliveChunkId(*ctx)
	id := lastAliveChunkId + 1
	numInsuranceUnbonded := len(state.InsuranceUnbonded)
	numUnpairedAliveChunk := len(unpairedAliveChunks) + numInsuranceUnbonded
	numPairedInsuranceBid := 0
	for _, elem := range rankedItems {
		switch v := elem.(type) {
		case *types.InsuranceBid:
			// pair qualified insurance bid
			var aliveChunk types.AliveChunk
			if numPairedInsuranceBid < numInsuranceUnbonded {
				unpaired := state.InsuranceUnbonded[numPairedInsuranceBid]
				aliveChunk = types.NewAliveChunk(id,
					// only token amount is used
					types.ChunkBondRequest{TokenAmount: unpaired.TokenAmount},
					*v,
				)
				delegationState.ChangeDelegation(unpaired.ValidatorAddress, v.ValidatorAddress, aliveChunk.TokenAmount)

				// TODO: accept insurance unbond request
				k.DeleteInsuranceUnbondRequest(*ctx, unpaired.Id)
				k.DeleteAliveChunk(*ctx, unpaired.Id)
				id++
			} else if numPairedInsuranceBid < numUnpairedAliveChunk {
				unpaired := unpairedAliveChunks[numPairedInsuranceBid-numInsuranceUnbonded]
				aliveChunk = types.NewAliveChunk(unpaired.Id,
					// only token amount is used
					types.ChunkBondRequest{TokenAmount: unpaired.TokenAmount},
					*v,
				)
				delegationState.ChangeDelegation(unpaired.ValidatorAddress, v.ValidatorAddress, aliveChunk.TokenAmount)
			} else {
				chunkBondRequest := state.ChunkBondRequests[numPairedInsuranceBid-numUnpairedAliveChunk]
				aliveChunk = types.NewAliveChunk(id, chunkBondRequest, *v)
				k.DeleteChunkBondRequest(*ctx, chunkBondRequest.Id)
				k.DeleteInsuranceBid(*ctx, v.Id)

				delegationState.ChangeDelegation("", v.ValidatorAddress, aliveChunk.TokenAmount)
				id++
			}
			k.DeleteInsuranceBid(*ctx, v.Id)
			k.SetAliveChunk(*ctx, aliveChunk)
			numPairedInsuranceBid++
		case *types.AliveChunk:
			if _, found := k.GetAliveChunk(*ctx, v.Id); !found {
				k.Logger(*ctx).Error("alive chunk has changed", "alive chunk id", v.Id)
				k.SetAliveChunk(*ctx, *v)
			}
		default:
			panic("items must be alive chunk or insurance bid")
		}
	}

	if numInsuranceUnbonded > numPairedInsuranceBid {
		numInsuranceUnbonded = numPairedInsuranceBid
	}
	state.InsuranceUnbonded = state.InsuranceUnbonded[numInsuranceUnbonded:]
	if id != lastAliveChunkId {
		k.SetLastAliveChunkId(*ctx, id)
	}
	return delegationState, nil
}

func (k *Keeper) resolveDelegationState(ctx *sdk.Context, delegationState types.DelegationState) error {
	for i := 0; i < len(delegationState.SortedValidators); i++ {
		var minVal, maxVal types.LiquidValidator
		min, max := sdk.ZeroInt(), sdk.ZeroInt()
		for _, val := range delegationState.SortedValidators {
			amount := delegationState.DelegationMap[val.OperatorAddress]
			if amount.LT(min) {
				min = amount
				minVal = val
			} else if amount.GT(max) {
				max = amount
				maxVal = val
			}
		}

		if min.IsZero() || max.IsZero() {
			break
		}
		tokenAmount := min.Abs()
		if tokenAmount.GT(max) {
			tokenAmount = max
		}
		if err := k.RedelegateTokenAmount(*ctx, minVal.OperatorAddress, maxVal.OperatorAddress, tokenAmount); err != nil {
			return err
		}
		delegationState.ChangeDelegation(maxVal.OperatorAddress, minVal.OperatorAddress, tokenAmount)
	}

	// resolve remaining
	for _, val := range delegationState.SortedValidators {
		tokenAmount := delegationState.DelegationMap[val.OperatorAddress]
		if tokenAmount.IsZero() {
			continue
		}
		if tokenAmount.IsPositive() {
			err := k.DelegateTokenAmount(*ctx, val.OperatorAddress, tokenAmount)
			if err != nil {
				return err
			}
			delegationState.ChangeDelegation("", val.OperatorAddress, tokenAmount)
		} else {
			panic("invalid scenario")
		}
	}
	return nil
}

func (k Keeper) PairChunkAndInsurance(ctx sdk.Context) error {
	state := types.State{
		InsuranceBids:           k.GetAllInsuranceBids(ctx),
		InsuranceUnbondRequests: k.GetAllInsuranceUnbondRequests(ctx),
		ChunkBondRequests:       k.GetAllChunkBondRequests(ctx),
		ChunkUnbondRequests:     k.GetAllChunkUnbondRequests(ctx),
		AliveChunks:             k.GetAllAliveChunks(ctx),
	}

	state = k.ResolveUnbondingQueues(ctx, state)

	// TODO: numAliveChunks is inaccurate (unranked alive pair should not be included)
	numAliveChunks := int64(len(state.AliveChunks) + len(state.ChunkBondRequests) +
		len(state.InsuranceUnbonded) - len(state.ChunkUnbondRequests))
	if numAliveChunks < 0 {
		panic("Chunk unbond request must be smaller than existing chunks")
	}
	intNumNewAliveChunks := sdk.NewInt(numAliveChunks)
	maxAliveChunks := k.GetParams(ctx).MaxAliveChunk
	intNumNewAliveChunks = sdk.MinInt(intNumNewAliveChunks, maxAliveChunks)
	numAliveChunks = intNumNewAliveChunks.Int64()

	ranked := k.rankAliveChunksAndInsuranceBid(&ctx, &state) // TODO: check. it is already sorted
	delegationState, err := k.unpairPairRankedAlivedChunks(&ctx, numAliveChunks, ranked, &state)
	if err != nil {
		return err
	}

	{
		for _, insuranceUnbonded := range state.InsuranceUnbonded {
			k.DeleteAliveChunk(ctx, insuranceUnbonded.Id)
			k.DeleteInsuranceUnbondRequest(ctx, insuranceUnbonded.Id)
			// TODO: check logic
			delegationState.ChangeDelegation("", insuranceUnbonded.ValidatorAddress, insuranceUnbonded.TokenAmount)
		}
		state.InsuranceUnbonded = types.InsuranceUnbondRequestedAliveChunks{}

		if err := k.resolveDelegationState(&ctx, delegationState); err != nil {
			return err
		}

		for _, chunkUnbonded := range state.ChunkUnbonded {
			if _, found := k.GetInsuranceUnbondRequest(ctx, chunkUnbonded.Id); found {
				// TODO: check spec, insurance unbond request is converted chunk unbond request
				//       due to lack of available chunk
				k.DeleteInsuranceUnbondRequest(ctx, chunkUnbonded.Id)
			}
			k.DeleteAliveChunk(ctx, chunkUnbonded.Id)
			if _, err := k.UndelegateTokenAmount(ctx,
				chunkUnbonded.Address,
				chunkUnbonded.ValidatorAddress,
				chunkUnbonded.TokenAmount); err != nil {
				return err
			}
		}
	}

	return nil
}
