package app

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/oracle/abci"
)

// PreBlocker is run before finalize block to update the aggregrate exchange rate votes on the oracle module
// that were verified by the vote etension handler so that the exchange rate votes are available during the
// entire block execution (from BeginBlock). It will execute the preblockers of the other modules set in
// SetOrderPreBlockers as well.
func (app *App) PreBlocker(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	if req == nil {
		err := fmt.Errorf("preblocker received a nil request")
		app.Logger().Error(err.Error())
		return nil, err
	}

	// execute preblockers of modules in OrderPreBlockers first.
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	paramsChanged := false
	for _, moduleName := range app.mm.OrderPreBlockers {
		if module, ok := app.mm.Modules[moduleName].(appmodule.HasPreBlocker); ok {
			rsp, err := module.PreBlock(ctx)
			if err != nil {
				return nil, err
			}
			if rsp.IsConsensusParamsChanged() {
				paramsChanged = true
			}
		}
	}

	res := &sdk.ResponsePreBlock{
		ConsensusParamsChanged: paramsChanged,
	}

	if len(req.Txs) == 0 {
		return res, nil
	}
	voteExtensionsEnabled := abci.VoteExtensionsEnabled(ctx)
	if voteExtensionsEnabled {
		var injectedVoteExtTx abci.AggregateExchangeRateVotes
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			app.Logger().Error("failed to decode injected vote extension tx", "err", err)
			return nil, err
		}

		// set oracle exchange rate votes using the passed in context, which will make
		// these votes available in the current block.
		for _, exchangeRateVote := range injectedVoteExtTx.ExchangeRateVotes {
			valAddr, err := sdk.ValAddressFromBech32(exchangeRateVote.Voter)
			if err != nil {
				app.Logger().Error("failed to get voter address", "err", err)
				continue
			}
			app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr, exchangeRateVote)
		}
	}

	app.Logger().Info(
		"oracle preblocker executed",
		"vote_extensions_enabled", voteExtensionsEnabled,
	)

	return res, nil
}
