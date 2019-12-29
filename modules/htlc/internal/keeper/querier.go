package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// NewQuerier creates a new HTLC Querier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryHTLC:
			return queryHTLC(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown HTLC query endpoint")
		}
	}
}

func queryHTLC(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryHTLCParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if len(params.HashLock) != types.HashLockLength {
		return nil, types.ErrInvalidHashLock(types.DefaultCodespace, fmt.Sprintf("the hash lock must be %d bytes long", types.HashLockLength))
	}

	htlc, err := keeper.GetHTLC(ctx, params.HashLock)
	if err != nil {
		return nil, err
	}

	bz, errRes := codec.MarshalJSONIndent(keeper.cdc, htlc)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}

	return bz, nil
}