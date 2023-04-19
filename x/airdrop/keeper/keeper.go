package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func (k Keeper) SetAirdropAccount(ctx sdk.Context, account types.AirdropAccount) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.AirdropAccountKey(account.OriginAddress),
		k.cdc.MustMarshal(&account),
	)
}

func (k Keeper) GetAirdropAccount(
	ctx sdk.Context, originAddress string,
) (account types.AirdropAccount, err error) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropAccountKey(originAddress))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshal(bz, &account)
	return
}