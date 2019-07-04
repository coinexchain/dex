package asset

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/tendermint/tendermint/libs/bech32"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	dex "github.com/coinexchain/dex/types"
)

var (
	SeparateKeyPrefix   = []byte{0x3A}
	TokenStoreKeyPrefix = []byte{0x01}
	WhitelistKeyPrefix  = []byte{0x02}
	ForbidAddrKeyPrefix = []byte{0x03}
)

// -----------------------------------------------------------------------------

// Keeper defines a module interface that keep token info.
type Keeper interface {
	TokenKeeper

	IssueToken(ctx sdk.Context, msg MsgIssueToken) sdk.Error
	TransferOwnership(ctx sdk.Context, msg MsgTransferOwnership) sdk.Error
	MintToken(ctx sdk.Context, msg MsgMintToken) sdk.Error
	BurnToken(ctx sdk.Context, msg MsgBurnToken) sdk.Error
	ForbidToken(ctx sdk.Context, msg MsgForbidToken) sdk.Error
	UnForbidToken(ctx sdk.Context, msg MsgUnForbidToken) sdk.Error
	AddTokenWhitelist(ctx sdk.Context, msg MsgAddTokenWhitelist) sdk.Error
	RemoveTokenWhitelist(ctx sdk.Context, msg MsgRemoveTokenWhitelist) sdk.Error
	ForbidAddress(ctx sdk.Context, msg MsgForbidAddr) sdk.Error
	UnForbidAddress(ctx sdk.Context, msg MsgUnForbidAddr) sdk.Error
	ModifyTokenURL(ctx sdk.Context, msg MsgModifyTokenURL) sdk.Error
	ModifyTokenDescription(ctx sdk.Context, msg MsgModifyTokenDescription) sdk.Error

	DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	AddToken(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SubtractToken(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SetParams(ctx sdk.Context, params Params)
	GetParams(ctx sdk.Context) (params Params)
}

var _ Keeper = (*BaseKeeper)(nil)

// BaseKeeper encodes/decodes tokens using the go-amino (binary) encoding/decoding library.
type BaseKeeper struct {
	BaseTokenKeeper

	// The codec codec for	binary encoding/decoding of token.
	cdc *codec.Codec
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	paramSubspace params.Subspace

	bkx ExpectedBankxKeeper
	sk  *staking.Keeper
}

// NewBaseKeeper returns a new BaseKeeper that uses go-amino to (binary) encode and decode concrete Token.
func NewBaseKeeper(cdc *codec.Codec, key sdk.StoreKey,
	paramStore params.Subspace, bkx ExpectedBankxKeeper, sk *staking.Keeper) BaseKeeper {
	return BaseKeeper{
		BaseTokenKeeper: NewBaseTokenKeeper(cdc, key),

		cdc:           cdc,
		key:           key,
		paramSubspace: paramStore.WithKeyTable(ParamKeyTable()),
		bkx:           bkx,
		sk:            sk,
	}
}

// DeductFee - deduct asset func fee like issueFee
func (keeper BaseKeeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {

	return keeper.bkx.DeductFee(ctx, addr, amt)
}

// AddToken - add token to addr when issue token and mint token etc.
func (keeper BaseKeeper) AddToken(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {

	return keeper.bkx.AddCoins(ctx, addr, amt)
}

// SubtractToken - sub token to addr when burn token etc.
func (keeper BaseKeeper) SubtractToken(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {

	return keeper.bkx.SubtractCoins(ctx, addr, amt)
}

// SetParams sets the asset module's parameters.
func (keeper BaseKeeper) SetParams(ctx sdk.Context, params Params) {
	keeper.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (keeper BaseKeeper) GetParams(ctx sdk.Context) (params Params) {
	keeper.paramSubspace.GetParamSet(ctx, &params)
	return
}

//IssueToken - new token and store it
func (keeper BaseKeeper) IssueToken(ctx sdk.Context, msg MsgIssueToken) sdk.Error {
	if keeper.IsTokenExists(ctx, msg.Symbol) {
		return ErrorDuplicateTokenSymbol(fmt.Sprintf("token symbol already exists in store"))
	}

	// only cet owner can issue reserved token
	if isReserved(msg.Symbol) && msg.Symbol != dex.CET {
		cetToken := keeper.GetToken(ctx, dex.CET)
		if cetToken == nil || !msg.Owner.Equals(cetToken.GetOwner()) {
			return ErrorInvalidTokenOwner("only coinex dex foundation can issue reserved symbol token, you can run \n" +
				"$ cetcli query asset reserved-symbol \n" +
				"to query reserved token symbol")
		}
	}

	token, err := NewToken(
		msg.Name,
		msg.Symbol,
		msg.TotalSupply,
		msg.Owner,
		msg.Mintable,
		msg.Burnable,
		msg.AddrForbiddable,
		msg.TokenForbiddable,
		msg.URL,
		msg.Description,
	)

	if err != nil {
		return err
	}

	return keeper.setToken(ctx, token)
}

//TransferOwnership - transfer token owner
func (keeper BaseKeeper) TransferOwnership(ctx sdk.Context, msg MsgTransferOwnership) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OriginalOwner)
	if err != nil {
		return err
	}

	if err := token.SetOwner(msg.NewOwner); err != nil {
		return ErrorInvalidTokenOwner("token new owner is invalid")
	}

	return keeper.setToken(ctx, token)
}

//MintToken - mint token
func (keeper BaseKeeper) MintToken(ctx sdk.Context, msg MsgMintToken) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetMintable() {
		return ErrorInvalidTokenMint(fmt.Sprintf("token %s do not support mint", msg.Symbol))
	}

	if err := token.SetTotalMint(msg.Amount + token.GetTotalMint()); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}

	if err := token.SetTotalSupply(msg.Amount + token.GetTotalSupply()); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	return keeper.setToken(ctx, token)
}

//BurnToken - burn token
func (keeper BaseKeeper) BurnToken(ctx sdk.Context, msg MsgBurnToken) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetBurnable() {
		return ErrorInvalidTokenBurn(fmt.Sprintf("token %s do not support burn", msg.Symbol))
	}

	if err := token.SetTotalBurn(msg.Amount + token.GetTotalBurn()); err != nil {
		return ErrorInvalidTokenBurn(err.Error())
	}

	if err := token.SetTotalSupply(token.GetTotalSupply() - msg.Amount); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	if token.GetSymbol() == dex.CET {
		updateBondPoolStatus(msg, keeper, ctx)
	}

	return keeper.setToken(ctx, token)
}

func updateBondPoolStatus(msg MsgBurnToken, keeper BaseKeeper, ctx sdk.Context) {
	decreaseNotBondedAmt := sdk.NewInt(msg.Amount).Neg()
	keeper.sk.InflateSupply(ctx, decreaseNotBondedAmt)
}

//ForbidToken - forbid token
func (keeper BaseKeeper) ForbidToken(ctx sdk.Context, msg MsgForbidToken) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token", msg.Symbol))
	}
	if token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(true)

	return keeper.setToken(ctx, token)
}

//UnForbidToken - unforbid token
func (keeper BaseKeeper) UnForbidToken(ctx sdk.Context, msg MsgUnForbidToken) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support unforbid token", msg.Symbol))
	}
	if !token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has not been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(false)

	return keeper.setToken(ctx, token)
}

//AddTokenWhitelist - add token forbidden whitelist
func (keeper BaseKeeper) AddTokenWhitelist(ctx sdk.Context, msg MsgAddTokenWhitelist) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token and add whitelist", msg.Symbol))
	}
	if err = keeper.addWhitelist(ctx, msg.Symbol, msg.Whitelist); err != nil {
		return ErrorInvalidTokenWhitelist(fmt.Sprintf("token whitelist is invalid"))
	}
	return nil
}

//RemoveTokenWhitelist - remove token forbidden whitelist
func (keeper BaseKeeper) RemoveTokenWhitelist(ctx sdk.Context, msg MsgRemoveTokenWhitelist) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token and remove whitelist", msg.Symbol))
	}
	if err = keeper.removeWhitelist(ctx, msg.Symbol, msg.Whitelist); err != nil {
		return ErrorInvalidTokenWhitelist(fmt.Sprintf("token whitelist is invalid"))
	}
	return nil
}

//ForbidAddress - add forbidden addresses
func (keeper BaseKeeper) ForbidAddress(ctx sdk.Context, msg MsgForbidAddr) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddr)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid address", msg.Symbol))
	}
	if err = keeper.addForbidAddress(ctx, msg.Symbol, msg.ForbidAddr); err != nil {
		return ErrorInvalidAddress(fmt.Sprintf("forbid addr is invalid"))
	}
	return nil
}

//UnForbidAddress - remove forbidden addresses
func (keeper BaseKeeper) UnForbidAddress(ctx sdk.Context, msg MsgUnForbidAddr) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddr)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support unforbid address", msg.Symbol))
	}
	if err = keeper.removeForbidAddress(ctx, msg.Symbol, msg.UnForbidAddr); err != nil {
		return ErrorInvalidAddress(fmt.Sprintf("unforbid addr is invalid"))
	}
	return nil
}

//ModifyTokenURL - modify token url property
func (keeper BaseKeeper) ModifyTokenURL(ctx sdk.Context, msg MsgModifyTokenURL) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if err := token.SetURL(msg.URL); err != nil {
		return ErrorInvalidTokenURL(err.Error())
	}

	return keeper.setToken(ctx, token)
}

//ModifyTokenURL - modify token url property
func (keeper BaseKeeper) ModifyTokenDescription(ctx sdk.Context, msg MsgModifyTokenDescription) sdk.Error {
	token, err := keeper.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if err := token.SetDescription(msg.Description); err != nil {
		return ErrorInvalidTokenDescription(err.Error())
	}

	return keeper.setToken(ctx, token)
}

func (keeper BaseKeeper) checkPrecondition(ctx sdk.Context, msg sdk.Msg, symbol string, owner sdk.AccAddress) (Token, sdk.Error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	token := keeper.GetToken(ctx, symbol)
	if token == nil {
		return nil, ErrorTokenNotFound(fmt.Sprintf("token %s not found", symbol))
	}

	if !token.GetOwner().Equals(owner) {
		return nil, ErrorInvalidTokenOwner("Only token owner can do this action")
	}

	return token, nil
}

func (keeper BaseKeeper) setToken(ctx sdk.Context, token Token) sdk.Error {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.key)

	bz, err := keeper.cdc.MarshalBinaryBare(token)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(TokenStoreKey(symbol), bz)
	return nil
}

func (keeper BaseKeeper) removeToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.key)
	store.Delete(TokenStoreKey(symbol))
}

func (keeper BaseKeeper) addWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.key)
	for _, addr := range whitelist {
		store.Set(PrefixAddrStoreKey(WhitelistKeyPrefix, symbol, addr), []byte{})
	}

	return nil
}

func (keeper BaseKeeper) removeWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.key)
	for _, addr := range whitelist {
		store.Delete(PrefixAddrStoreKey(WhitelistKeyPrefix, symbol, addr))
	}

	return nil
}

func (keeper BaseKeeper) addForbidAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.key)
	for _, addr := range addresses {
		store.Set(PrefixAddrStoreKey(ForbidAddrKeyPrefix, symbol, addr), []byte{})
	}

	return nil
}

func (keeper BaseKeeper) removeForbidAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.key)
	for _, addr := range addresses {
		store.Delete(PrefixAddrStoreKey(ForbidAddrKeyPrefix, symbol, addr))
	}

	return nil
}

// -----------------------------------------------------------------------------

// TokenKeeper defines a module interface that facilitates read only access to token info.
type TokenKeeper interface {
	ViewKeeper

	IsTokenForbidden(ctx sdk.Context, symbol string) bool
	IsTokenExists(ctx sdk.Context, symbol string) bool
	IsTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool
	IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool
}

var _ TokenKeeper = (*BaseTokenKeeper)(nil)

// BaseTokenKeeper implements a read only keeper implementation of TokenKeeper.
type BaseTokenKeeper struct {
	BaseViewKeeper

	// The codec codec for	binary encoding/decoding of token.
	cdc *codec.Codec
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
}

// NewBaseTokenKeeper returns a new NewBaseTokenKeeper that uses go-amino to (binary) encode and decode concrete Token.
func NewBaseTokenKeeper(cdc *codec.Codec, key sdk.StoreKey) BaseTokenKeeper {
	return BaseTokenKeeper{
		BaseViewKeeper: NewBaseViewKeeper(cdc, key),
		cdc:            cdc,
		key:            key,
	}
}

//IsTokenForbidden - check whether coin issuer has forbidden "denom"
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
	if token == nil {
		return true
	}

	if keeper.hasAddrKey(ctx, ForbidAddrKeyPrefix, symbol, addr) {
		return true
	}

	if !token.GetIsForbidden() {
		return false
	}

	if keeper.hasAddrKey(ctx, WhitelistKeyPrefix, symbol, addr) {
		return false
	}

	if token.GetOwner().Equals(addr) {
		return false
	}

	return true
}

// hasAddrKey - KV store KEY: prefix | symbol: | AccAddress
func (keeper BaseTokenKeeper) hasAddrKey(ctx sdk.Context, prefix []byte, symbol string, addr sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.key)
	key := PrefixAddrStoreKey(prefix, symbol, addr)
	return store.Has(key)
}
func (keeper BaseTokenKeeper) setAddrKey(ctx sdk.Context, prefix []byte, addr string) error {
	store := ctx.KVStore(keeper.key)
	index := strings.Index(addr, string(SeparateKeyPrefix))

	accBech32, err := sdk.AccAddressFromBech32(string([]byte(addr)[index+1:]))
	if err != nil {
		return err
	}
	key := PrefixAddrStoreKey(prefix, string([]byte(addr)[:index]), accBech32)
	store.Set(key, []byte{})

	return nil
}

// -----------------------------------------------------------------------------

// ViewKeeper defines a module interface that facilitates read only access to token store info.
type ViewKeeper interface {
	GetToken(ctx sdk.Context, symbol string) Token
	GetAllTokens(ctx sdk.Context) []Token
	GetWhitelist(ctx sdk.Context, symbol string) []sdk.AccAddress
	GetForbiddenList(ctx sdk.Context, symbol string) []sdk.AccAddress
	GetAllAddrKeys(ctx sdk.Context, prefix []byte) []string
	GetReservedSymbols() []string
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	// The codec codec for	binary encoding/decoding of token.
	cdc *codec.Codec
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
}

// BaseViewKeeper returns a new BaseViewKeeper that uses go-amino to (binary) encode and decode concrete Token.
func NewBaseViewKeeper(cdc *codec.Codec, key sdk.StoreKey) BaseViewKeeper {
	return BaseViewKeeper{
		cdc: cdc,
		key: key,
	}
}

// GetToken - return token by symbol
func (keeper BaseViewKeeper) GetToken(ctx sdk.Context, symbol string) Token {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(TokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	return keeper.decodeToken(bz)
}

// GetAllTokens - returns all tokens.
func (keeper BaseViewKeeper) GetAllTokens(ctx sdk.Context) []Token {
	tokens := make([]Token, 0)

	keeper.iterateTokenValue(ctx, func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return tokens
}

// GetWhitelist - returns whitelist.
func (keeper BaseViewKeeper) GetWhitelist(ctx sdk.Context, symbol string) []sdk.AccAddress {
	whitelist := make([]sdk.AccAddress, 0)
	keyPrefix := append(append(WhitelistKeyPrefix, symbol...), SeparateKeyPrefix...)

	keeper.iterateAddrKey(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[len(WhitelistKeyPrefix)+len(symbol)+len(SeparateKeyPrefix):]
		whitelist = append(whitelist, addr)
		return false
	})

	return whitelist
}

// GetForbidAddr - returns all forbidden addr list.
func (keeper BaseViewKeeper) GetForbiddenList(ctx sdk.Context, symbol string) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0)
	keyPrefix := append(append(ForbidAddrKeyPrefix, symbol...), SeparateKeyPrefix...)

	keeper.iterateAddrKey(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[len(ForbidAddrKeyPrefix)+len(symbol)+len(SeparateKeyPrefix):]
		addresses = append(addresses, addr)
		return false
	})

	return addresses
}

// GetAllAddrKeys return []KEY symbol: | addr . get all whitelists or forbidden addresses string to genesis.json
func (keeper BaseViewKeeper) GetAllAddrKeys(ctx sdk.Context, prefix []byte) []string {
	res := make([]string, 0)
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	keeper.iterateAddrKey(ctx, prefix, func(key []byte) (stop bool) {
		i := bytes.Index(key, SeparateKeyPrefix) + len(SeparateKeyPrefix)
		bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, key[i:])
		if err != nil {
			panic(err)
		}
		s := string(key[len(prefix):i]) + bech32Addr
		res = append(res, s)
		return false
	})

	return res
}

// GetReservedSymbols - get all reserved symbols
func (keeper BaseViewKeeper) GetReservedSymbols() []string {
	return reserved
}

func (keeper BaseViewKeeper) iterateTokenValue(ctx sdk.Context, process func(Token) (stop bool)) {
	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		acc := keeper.decodeToken(iter.Value())
		if process(acc) {
			return
		}
		iter.Next()
	}
}
func (keeper BaseViewKeeper) iterateAddrKey(ctx sdk.Context, prefix []byte, process func(key []byte) (stop bool)) {
	store := ctx.KVStore(keeper.key)
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
func (keeper BaseViewKeeper) decodeToken(bz []byte) (token Token) {
	if err := keeper.cdc.UnmarshalBinaryBare(bz, &token); err != nil {
		panic(err)
	}
	return
}

// -----------------------------------------------------------------------------

// TokenStoreKey turn token symbol to KEY prefix | symbol .
func TokenStoreKey(symbol string) []byte {
	return append(TokenStoreKeyPrefix, []byte(symbol)...)
}

// PrefixAddrStoreKey - new KEY prefix | Symbol: | AccAddress
func PrefixAddrStoreKey(prefix []byte, symbol string, addr sdk.AccAddress) []byte {
	return append(append(append(prefix, symbol...), SeparateKeyPrefix...), addr...)
}
