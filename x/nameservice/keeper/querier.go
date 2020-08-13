package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Hansol-git/nameservice/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the nameservice Querier
const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

// path : [queryType, name, ]

// NewQuerier creates a new querier for nameservice clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, k)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, k)
		case QueryNames:
			return queryNames(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}

// queryResolve - func for queryResolve
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	value := k.ResolveName(ctx, path[0])

	if value == "" {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not resolve name")
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.QueryResResolve{Value: value})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryWhois - func for queryWhois
func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	whois := k.GetWhois(ctx, path[0])

	res, err := codec.MarshalJSONIndent(k.cdc, whois)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryNames - func for queryNames
func queryNames(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var nameList types.QueryResNames

	iterator := k.GetNamesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		nameList = append(nameList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, nameList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
