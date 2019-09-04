package keepers

import (
	"bytes"
	"errors"
	"strings"

	"github.com/tendermint/tendermint/libs/bech32"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
)

// -----------------------------------------------------------------------------

// Keeper defines a module interface that keep token info.
type Keeper interface {
	TokenKeeper

	IssueToken(ctx sdk.Context, name string, symbol string, totalSupply sdk.Int, owner sdk.AccAddress,
		mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool,
		url string, description string, identity string) sdk.Error
	TransferOwnership(ctx sdk.Context, symbol string, originalOwner sdk.AccAddress, newOwner sdk.AccAddress) sdk.Error
	MintToken(ctx sdk.Context, symbol string, owner sdk.AccAddress, amount sdk.Int) sdk.Error
	BurnToken(ctx sdk.Context, symbol string, owner sdk.AccAddress, amount sdk.Int) sdk.Error
	ForbidToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) sdk.Error
	UnForbidToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) sdk.Error
	AddTokenWhitelist(ctx sdk.Context, symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) sdk.Error
	RemoveTokenWhitelist(ctx sdk.Context, symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) sdk.Error
	ForbidAddress(ctx sdk.Context, symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) sdk.Error
	UnForbidAddress(ctx sdk.Context, symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) sdk.Error
	ModifyTokenInfo(ctx sdk.Context, symbol string, owner sdk.AccAddress, url string, description string) sdk.Error

	SetParams(ctx sdk.Context, params types.Params)
	GetParams(ctx sdk.Context) (params types.Params)
}

var _ Keeper = (*BaseKeeper)(nil)

// BaseKeeper encodes/decodes tokens using the go-amino (binary) encoding/decoding library.
type BaseKeeper struct {
	BaseTokenKeeper

	// The codec codec for	binary encoding/decoding of token.
	cdc *codec.Codec
	// The (unexposed) key used to access the store from the Context.
	storeKey sdk.StoreKey

	paramSubspace params.Subspace

	bkx types.ExpectedBankxKeeper
	sk  types.ExpectedSupplyKeeper
}

// NewBaseKeeper returns a new BaseKeeper that uses go-amino to (binary) encode and decode concrete Token.
func NewBaseKeeper(cdc *codec.Codec, key sdk.StoreKey,
	paramStore params.Subspace, bkx types.ExpectedBankxKeeper, sk supply.Keeper) BaseKeeper {
	return BaseKeeper{
		BaseTokenKeeper: NewBaseTokenKeeper(cdc, key),

		cdc:           cdc,
		storeKey:      key,
		paramSubspace: paramStore.WithKeyTable(ParamKeyTable()),
		bkx:           bkx,
		sk:            sk,
	}
}

// IssueToken - new token and store it
func (keeper BaseKeeper) IssueToken(ctx sdk.Context, name string, symbol string, totalSupply sdk.Int, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool,
	url string, description string, identity string) sdk.Error {
	if keeper.bkx.BlacklistedAddr(owner) {
		return types.ErrAccInBlackList(owner)
	}

	if keeper.IsTokenExists(ctx, symbol) {
		return types.ErrDuplicateTokenSymbol(symbol)
	}

	var cetToken types.Token
	// only cet owner can issue reserved token
	if types.IsReservedSymbol(symbol) && symbol != dex.CET {
		cetToken = keeper.GetToken(ctx, dex.CET)
		if cetToken == nil || !owner.Equals(cetToken.GetOwner()) {
			return types.ErrInvalidIssueOwner()
		}
	}

	// only cet owner can issue .suffix token
	if types.IsSuffixSymbol(symbol) {
		cetToken = keeper.GetToken(ctx, dex.CET)
		if cetToken == nil || !owner.Equals(cetToken.GetOwner()) {
			return types.ErrInvalidTokenSymbol(symbol)
		}
	}

	token, err := types.NewToken(
		name,
		symbol,
		totalSupply,
		owner,
		mintable,
		burnable,
		addrForbiddable,
		tokenForbiddable,
		url,
		description,
		identity,
	)

	if err != nil {
		return err
	}

	if err := keeper.SetToken(ctx, token); err != nil {
		return err
	}

	return keeper.sk.MintCoins(ctx, types.ModuleName, types.NewTokenCoins(symbol, totalSupply))
}

// TransferOwnership - transfer token owner
func (keeper BaseKeeper) TransferOwnership(ctx sdk.Context, symbol string, originalOwner sdk.AccAddress, newOwner sdk.AccAddress) sdk.Error {
	if keeper.bkx.BlacklistedAddr(newOwner) {
		return types.ErrAccInBlackList(newOwner)
	}
	token, err := keeper.checkPrecondition(ctx, symbol, originalOwner)
	if err != nil {
		return err
	}

	if err := token.SetOwner(newOwner); err != nil {
		return err
	}

	return keeper.SetToken(ctx, token)
}

// MintToken - mint token
func (keeper BaseKeeper) MintToken(ctx sdk.Context, symbol string, owner sdk.AccAddress, amount sdk.Int) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetMintable() {
		return types.ErrTokenMintNotSupported(symbol)
	}

	if err := token.SetTotalMint(token.GetTotalMint().Add(amount)); err != nil {
		return err
	}

	if err := token.SetTotalSupply(token.GetTotalSupply().Add(amount)); err != nil {
		return err
	}

	if err := keeper.SetToken(ctx, token); err != nil {
		return err
	}

	return keeper.sk.MintCoins(ctx, types.ModuleName, types.NewTokenCoins(symbol, amount))
}

// BurnToken - burn token
func (keeper BaseKeeper) BurnToken(ctx sdk.Context, symbol string, owner sdk.AccAddress, amount sdk.Int) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetBurnable() {
		return types.ErrTokenBurnNotSupported(symbol)
	}

	if err := token.SetTotalBurn(token.GetTotalBurn().Add(amount)); err != nil {
		return err
	}

	if err := token.SetTotalSupply(token.GetTotalSupply().Sub(amount)); err != nil {
		return err
	}

	if err := keeper.SetToken(ctx, token); err != nil {
		return err
	}

	return keeper.sk.BurnCoins(ctx, types.ModuleName, types.NewTokenCoins(symbol, amount))

}

// ForbidToken - forbid token
func (keeper BaseKeeper) ForbidToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return types.ErrTokenForbiddenNotSupported(symbol)
	}
	if token.GetIsForbidden() {
		return types.ErrInvalidTokenForbidden(symbol)
	}
	token.SetIsForbidden(true)

	return keeper.SetToken(ctx, token)
}

// UnForbidToken - unforbid token
func (keeper BaseKeeper) UnForbidToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return types.ErrTokenForbiddenNotSupported(symbol)
	}
	if !token.GetIsForbidden() {
		return types.ErrInvalidTokenUnForbidden(symbol)
	}
	token.SetIsForbidden(false)

	return keeper.SetToken(ctx, token)
}

// AddTokenWhitelist - add token forbidden whitelist
func (keeper BaseKeeper) AddTokenWhitelist(ctx sdk.Context, symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return types.ErrTokenForbiddenNotSupported(symbol)
	}
	if err = keeper.addWhitelist(ctx, symbol, whitelist); err != nil {
		return types.ErrInvalidTokenWhitelist()
	}
	return nil
}

// RemoveTokenWhitelist - remove token forbidden whitelist
func (keeper BaseKeeper) RemoveTokenWhitelist(ctx sdk.Context, symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return types.ErrTokenForbiddenNotSupported(symbol)
	}
	if err = keeper.removeWhitelist(ctx, symbol, whitelist); err != nil {
		return types.ErrInvalidTokenWhitelist()
	}
	return nil
}

// ForbidAddress - add forbidden addresses
func (keeper BaseKeeper) ForbidAddress(ctx sdk.Context, symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return types.ErrAddressForbiddenNotSupported(symbol)
	}
	if err = keeper.addForbiddenAddress(ctx, symbol, addresses); err != nil {
		return types.ErrInvalidForbiddenAddress()
	}
	return nil
}

// UnForbidAddress - remove forbidden addresses
func (keeper BaseKeeper) UnForbidAddress(ctx sdk.Context, symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return types.ErrAddressForbiddenNotSupported(symbol)
	}
	if err = keeper.removeForbiddenAddress(ctx, symbol, addresses); err != nil {
		return types.ErrInvalidForbiddenAddress()
	}
	return nil
}

// ModifyTokenInfo - modify token info property
func (keeper BaseKeeper) ModifyTokenInfo(ctx sdk.Context, symbol string, owner sdk.AccAddress, url string, description string) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, symbol, owner)
	if err != nil {
		return err
	}

	if url != types.DoNotModifyTokenInfo {
		if err := token.SetURL(url); err != nil {
			return err
		}
	}

	if description != types.DoNotModifyTokenInfo {
		if err := token.SetDescription(description); err != nil {
			return err
		}
	}

	return keeper.SetToken(ctx, token)
}

func (keeper BaseKeeper) SendCoinsFromAssetModuleToAccount(ctx sdk.Context, addresses sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return keeper.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addresses, amt)
}

func (keeper BaseKeeper) SendCoinsFromAccountToAssetModule(ctx sdk.Context, addresses sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return keeper.sk.SendCoinsFromAccountToModule(ctx, addresses, types.ModuleName, amt)
}

// DeductIssueFee - deduct issue token fee
func (keeper BaseKeeper) DeductIssueFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return keeper.bkx.DeductFee(ctx, addr, amt)
}

// AddToken - used for unit test
func (keeper BaseKeeper) AddToken(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return keeper.bkx.AddCoins(ctx, addr, amt)
}

// GetAccTotalToken - used for unit test
func (keeper BaseKeeper) GetAccTotalToken(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return keeper.bkx.GetTotalCoins(ctx, addr)
}

func (keeper BaseKeeper) checkPrecondition(ctx sdk.Context, symbol string, owner sdk.AccAddress) (types.Token, sdk.Error) {
	token := keeper.GetToken(ctx, symbol)
	if token == nil {
		return nil, types.ErrTokenNotFound(symbol)
	}

	if !token.GetOwner().Equals(owner) {
		return nil, types.ErrNeedTokenOwner(token.GetOwner())
	}

	return token, nil
}

func (keeper BaseKeeper) RemoveToken(ctx sdk.Context, token types.Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.GetTokenStoreKey(symbol))
}

func (keeper BaseKeeper) addWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	for _, addr := range whitelist {
		store.Set(types.GetWhitelistStoreKey(symbol, addr), []byte{})
	}

	return nil
}

func (keeper BaseKeeper) removeWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	for _, addr := range whitelist {
		store.Delete(types.GetWhitelistStoreKey(symbol, addr))
	}

	return nil
}

func (keeper BaseKeeper) addForbiddenAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	for _, addr := range addresses {
		store.Set(types.GetForbiddenAddrStoreKey(symbol, addr), []byte{})
	}

	return nil
}

func (keeper BaseKeeper) removeForbiddenAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	for _, addr := range addresses {
		store.Delete(types.GetForbiddenAddrStoreKey(symbol, addr))
	}

	return nil
}

// -----------------------------------------------------------------------------

// TokenKeeper defines a module interface that facilitates read only access to token store info.
type TokenKeeper interface {
	GetToken(ctx sdk.Context, symbol string) types.Token
	GetAllTokens(ctx sdk.Context) []types.Token
	GetWhitelist(ctx sdk.Context, symbol string) []sdk.AccAddress
	GetForbiddenAddresses(ctx sdk.Context, symbol string) []sdk.AccAddress

	IsTokenForbidden(ctx sdk.Context, symbol string) bool
	IsTokenExists(ctx sdk.Context, symbol string) bool
	IsTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool
	IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool
	UpdateTokenSendLock(ctx sdk.Context, symbol string, amount sdk.Int, lock bool) sdk.Error
}

var _ TokenKeeper = (*BaseTokenKeeper)(nil)

// BaseTokenKeeper implements a read only keeper implementation of TokenKeeper.
type BaseTokenKeeper struct {
	// The codec codec for	binary encoding/decoding of token.
	cdc *codec.Codec
	// The (unexposed) key used to access the store from the Context.
	storeKey sdk.StoreKey
}

// BaseTokenKeeper returns a new BaseTokenKeeper that uses go-amino to (binary) encode and decode concrete Token.
func NewBaseTokenKeeper(cdc *codec.Codec, key sdk.StoreKey) BaseTokenKeeper {
	return BaseTokenKeeper{
		cdc:      cdc,
		storeKey: key,
	}
}

// GetToken - return token by symbol
func (keeper BaseTokenKeeper) GetToken(ctx sdk.Context, symbol string) types.Token {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.GetTokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	return types.MustUnmarshalToken(keeper.cdc, bz)
}

// GetAllTokens - returns all tokens.
func (keeper BaseTokenKeeper) GetAllTokens(ctx sdk.Context) []types.Token {
	tokens := make([]types.Token, 0)
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TokenKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		token := types.MustUnmarshalToken(keeper.cdc, iterator.Value())
		tokens = append(tokens, token)
	}
	return tokens
}

// GetWhitelist - returns whitelist.
func (keeper BaseTokenKeeper) GetWhitelist(ctx sdk.Context, symbol string) []sdk.AccAddress {
	whitelist := make([]sdk.AccAddress, 0)
	keyPrefix := types.GetWhitelistKeyPrefix(symbol)

	keeper.iterateAddrKey(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[types.GetWhitelistKeyPrefixLength(symbol):]
		whitelist = append(whitelist, addr)
		return false
	})

	return whitelist
}

// GetForbiddenAddresses - returns all forbidden addr
func (keeper BaseTokenKeeper) GetForbiddenAddresses(ctx sdk.Context, symbol string) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0)
	keyPrefix := types.GetForbiddenAddrKeyPrefix(symbol)

	keeper.iterateAddrKey(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[types.GetForbiddenAddrKeyPrefixLength(symbol):]
		addresses = append(addresses, addr)
		return false
	})

	return addresses
}

//IsTokenForbidden - check whether coin issuer has forbidden "symbol"
func (keeper BaseTokenKeeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	token := keeper.GetToken(ctx, symbol)
	if token != nil {
		return token.GetIsForbidden()
	}

	return true
}

// IsTokenExists - check whether there is a coin named "symbol"
func (keeper BaseTokenKeeper) IsTokenExists(ctx sdk.Context, symbol string) bool {
	return keeper.GetToken(ctx, symbol) != nil
}

// IsTokenIssuer - check whether addr is a token issuer
func (keeper BaseTokenKeeper) IsTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	if addr.Empty() {
		return false
	}

	token := keeper.GetToken(ctx, symbol)
	return token != nil && token.GetOwner().Equals(addr)
}

// IsForbiddenByTokenIssuer - check whether addr is forbid by token issuer
func (keeper BaseTokenKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	token := keeper.GetToken(ctx, symbol)
	store := ctx.KVStore(keeper.storeKey)
	if token == nil {
		return true
	}

	if store.Has(types.GetForbiddenAddrStoreKey(symbol, addr)) {
		return true
	}

	if !token.GetIsForbidden() {
		return false
	}

	if store.Has(types.GetWhitelistStoreKey(symbol, addr)) {
		return false
	}

	if token.GetOwner().Equals(addr) {
		return false
	}

	return true
}

// UpdateTokenSendLock - set token SendLock amount
func (keeper BaseTokenKeeper) UpdateTokenSendLock(ctx sdk.Context, symbol string, amount sdk.Int, lock bool) sdk.Error {
	token := keeper.GetToken(ctx, symbol)
	if token == nil {
		return types.ErrTokenNotFound(symbol)
	}
	if lock {
		if err := token.SetSendLock(token.GetSendLock().Add(amount)); err != nil {
			return err
		}
	} else {
		if err := token.SetSendLock(token.GetSendLock().Sub(amount)); err != nil {
			return err
		}
	}

	if err := keeper.SetToken(ctx, token); err != nil {
		return err
	}
	return nil

}

// SetToken - set token to store
func (keeper BaseTokenKeeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.storeKey)

	bz, err := keeper.cdc.MarshalBinaryBare(token)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(types.GetTokenStoreKey(symbol), bz)
	return nil
}

// ImportGenesisAddrKeys - import all whitelists or forbidden addresses string from genesis.json
func (keeper BaseTokenKeeper) ImportGenesisAddrKeys(ctx sdk.Context, prefix []byte, addr string) error {
	store := ctx.KVStore(keeper.storeKey)

	// symbol | : | address
	split := strings.SplitAfterN(addr, string(types.SeparateKey), 2)
	if len(split) != 2 {
		return errors.New("Genesis Address Err ")
	}
	addrBech32, err := sdk.AccAddressFromBech32(split[1])
	if err != nil {
		return err
	}
	key := append(append(prefix, split[0]...), addrBech32...)
	store.Set(key, []byte{})

	return nil
}

// ExportGenesisAddrKeys - get all whitelists or forbidden addresses string to genesis.json
func (keeper BaseTokenKeeper) ExportGenesisAddrKeys(ctx sdk.Context, prefix []byte) (res []string) {
	bech32AccountAddrPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	keeper.iterateAddrKey(ctx, prefix, func(key []byte) (stop bool) {

		// prefix | symbol | : | address
		split := bytes.SplitAfterN(key, types.SeparateKey, 2)
		if len(split) != 2 {
			panic(errors.New("Genesis Addr Err "))
		}
		addrBech32, err := bech32.ConvertAndEncode(bech32AccountAddrPrefix, split[1])
		if err != nil {
			panic(err)
		}
		s := string(split[0][len(prefix):]) + addrBech32
		res = append(res, s)
		return false
	})

	return res
}

func (keeper BaseTokenKeeper) iterateAddrKey(ctx sdk.Context, prefix []byte, process func(key []byte) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		key := iter.Key()
		if process(key) {
			return
		}
		iter.Next()
	}
}
