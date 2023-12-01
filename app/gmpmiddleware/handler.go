package gmpmiddleware

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type GmpKeeper interface {
	RelayPrice(
		goCtx context.Context,
		msg *types.MsgRelayPrice,
	) (*types.MsgRelayPriceResponse, error)
	GetParams(ctx sdk.Context) (params types.Params)
}

type GmpHandler struct {
	gmp GmpKeeper
}

func NewGmpHandler(k GmpKeeper) *GmpHandler {
	return &GmpHandler{
		gmp: k,
	}
}

// HandleGeneralMessage takes the receiving message from axelar,
// and sends it along to the GMP module.
func (h GmpHandler) HandleGeneralMessage(
	ctx sdk.Context,
	srcChain,
	srcAddress string,
	receiver string,
	payload []byte,
	sender string,
	channel string,
) error {
	ctx.Logger().Info("HandleGeneralMessage called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"receiver", receiver,
		"payload", payload,
		"module", "x/gmp-middleware",
	)

	err := verifyParams(h.gmp.GetParams(ctx), sender, channel)
	if err != nil {
		return err
	}
	msg, err := types.NewGmpDecoder(payload)
	if err != nil {
		return err
	}
	tx := &types.MsgRelayPrice{
		Relayer:               srcAddress,
		DestinationChain:      srcChain,
		ClientContractAddress: msg.ContractAddress.Hex(),
		OjoContractAddress:    srcAddress,
		Denoms:                msg.GetDenoms(),
		CommandSelector:       msg.CommandSelector[:],
		CommandParams:         msg.CommandParams,
		Timestamp:             msg.Timestamp.Int64(),
	}
	err = tx.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = h.gmp.RelayPrice(ctx, tx)
	return err
}

// HandleGeneralMessage takes the receiving message from axelar,
// and sends it along to the GMP module.
func (h GmpHandler) HandleGeneralMessageWithToken(
	ctx sdk.Context,
	srcChain,
	srcAddress string,
	receiver string,
	payload []byte,
	sender string,
	channel string,
	coin sdk.Coin,
) error {
	ctx.Logger().Info("HandleGeneralMessageWithToken called",
		"srcChain", srcChain,
		"srcAddress", srcAddress, // this is the Ojo contract address
		"receiver", receiver,
		"payload", payload,
		"coin", coin,
	)

	err := verifyParams(h.gmp.GetParams(ctx), sender, channel)
	if err != nil {
		return err
	}
	msg, err := types.NewGmpDecoder(payload)
	if err != nil {
		return err
	}
	tx := &types.MsgRelayPrice{
		Relayer:               srcAddress,
		DestinationChain:      srcChain,
		ClientContractAddress: msg.ContractAddress.Hex(),
		OjoContractAddress:    srcAddress,
		Denoms:                msg.GetDenoms(),
		CommandSelector:       msg.CommandSelector[:],
		CommandParams:         msg.CommandParams,
		Timestamp:             msg.Timestamp.Int64(),
		Token:                 coin,
	}
	err = tx.ValidateBasic()
	if err != nil {
		return err
	}
	_, err = h.gmp.RelayPrice(ctx, tx)
	return err
}